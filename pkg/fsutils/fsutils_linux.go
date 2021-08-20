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
