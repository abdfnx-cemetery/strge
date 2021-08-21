// +build windows

package constants

import (
	"fmt"
	"path/filepath"
	"strings"
)

func CheckSystemDriveAndRemoveDriveLetter(path string) (string, error) {
	if len(path) == 2 && string(path[1]) == ":" {
		return "", fmt.Errorf("No relative path specified in %q", path)
	}

	if !filepath.IsAbs(path) || len(path) < 2 {
		return filepath.FromSlash(path), nil
	}

	if string(path[1]) == ":" && !strings.EqualFold(string(path[0]), "c") {
		return "", fmt.Errorf("The specified path is not on the system drive (C:)")
	}

	return filepath.FromSlash(path[2:]), nil
}
