package context

import (
	"fmt"
	"syscall"
)

// chrootOrChdir() is either a chdir() to the specified path, or a chroot() to the
// specified path followed by chdir() to the new root directory
func chrootOrChdir(path string) error {
	if err := syscall.Chdir(path); err != nil {
		return fmt.Errorf("error changing to %q: %v", path, err)
	}

	return nil
}
