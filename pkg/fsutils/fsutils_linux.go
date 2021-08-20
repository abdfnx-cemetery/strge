// +build linux

package fsutils

import (
	"fmt"
	"io/ioutil"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func locateDummyIfEmpty(path string) (string, error) {
	children, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	if len(children) != 0 {
		return "", nil
	}

	dummyFile, err := ioutil.TempFile(path, "fsutils-dummy")
	if err != nil {
		return "", err
	}

	name := dummyFile.Name()
	err = dummyFile.Close()
	return name, err
}

// SupportsDType returns whether the filesystem mounted on path supports d_type
func SupportsDType(path string) (bool, error) {
	// locate dummy so that we have at least one dirent
	dummy, err := locateDummyIfEmpty(path)
	if err != nil {
		return false, err
	}
	if dummy != "" {
		defer os.Remove(dummy)
	}

	visited := 0
	supportsDType := true
	fn := func(ent *unix.Dirent) bool {
		visited++
		if ent.Type == unix.DT_UNKNOWN {
			supportsDType = false
			// stop iteration
			return true
		}
		// continue iteration
		return false
	}

	if err = iterateReadDir(path, fn); err != nil {
		return false, err
	}

	if visited == 0 {
		return false, fmt.Errorf("did not hit any dirent during iteration %s", path)
	}

	return supportsDType, nil
}
