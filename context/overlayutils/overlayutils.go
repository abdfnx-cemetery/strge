// +build linux

package overlayutils

import (
	"fmt"

	"github.com/gepis/strge/context"
	"github.com/pkg/errors"
)

// ErrDTypeNotSupported denotes that the backing filesystem doesn't support d_type.
func ErrDTypeNotSupported(driver, backingFs string) error {
	msg := fmt.Sprintf("%s: the backing %s filesystem is formatted without d_type support, which leads to incorrect behavior.", driver, backingFs)
	if backingFs == "xfs" {
		msg += " Reformat the filesystem with ftype=1 to enable d_type support."
	}
	msg += " Running without d_type is not supported."
	return errors.Wrap(context.ErrNotSupported, msg)
}
