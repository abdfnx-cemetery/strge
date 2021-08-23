// +build !go1.16

package overlay

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gepis/strge/pkg/archive"
	"github.com/gepis/strge/pkg/constants"
)

func scanForMountProgramIndicators(home string) (detected bool, err error) {
	err = filepath.Walk(home, func(path string, info os.FileInfo, err error) error {
		if detected {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		basename := filepath.Base(path)
		if strings.HasPrefix(basename, archive.WhiteoutPrefix) {
			detected = true
			return filepath.SkipDir
		}
		if info.IsDir() {
			xattrs, err := constants.Llistxattr(path)
			if err != nil {
				return err
			}
			for _, xattr := range xattrs {
				if strings.HasPrefix(xattr, "user.fuseoverlayfs.") || strings.HasPrefix(xattr, "user.containers.") {
					detected = true
					return filepath.SkipDir
				}
			}
		}
		return nil
	})
	return detected, err
}
