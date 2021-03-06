// +build !linux

package chunked

import (
	"context"

	storage "github.com/gepis/strge"
	cntx "github.com/gepis/strge/context"
	"github.com/pkg/errors"
)

// GetDiffer returns a differ than can be used with ApplyDiffWithDiffer.
func GetDiffer(ctx context.Context, store storage.Store, blobSize int64, annotations map[string]string, iss ImageSourceSeekable) (cntx.Differ, error) {
	return nil, errors.New("format not supported on this architecture")
}
