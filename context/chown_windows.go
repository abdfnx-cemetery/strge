// +build windows

package context

import (
	"os"
	"syscall"

	"github.com/gepis/strge/pkg/idtools"
)

func platformLChown(path string, info os.FileInfo, toHost, toContainer *idtools.IDMappings) error {
	return &os.PathError{"lchown", path, syscall.EWINDOWS}
}
