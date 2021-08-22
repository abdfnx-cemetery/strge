package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gepis/strge/pkg/idtools"
	"github.com/gepis/strge/pkg/ioutils"
	"github.com/gepis/strge/pkg/stringid"
	"github.com/gepis/strge/pkg/truncindex"
	digest "github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

type Container struct {
	ID string `json:"id"`
	Names []string `json:"names,omitempty"`
	ImageID string `json:"image"`
	LayerID string `json:"layer"`
	Metadata string `json:"metadata,omitempty"`
	BigDataNames []string `json:"big-data-names,omitempty"`
	BigDataSizes map[string]int64 `json:"big-data-sizes,omitempty"`
	BigDataDigests map[string]digest.Digest `json:"big-data-digests,omitempty"`
	Created time.Time `json:"created,omitempty"`
	UIDMap []idtools.IDMap `json:"uidmap,omitempty"`
	GIDMap []idtools.IDMap `json:"gidmap,omitempty"`
	Flags map[string]interface{} `json:"flags,omitempty"`
}

type ContainerStore interface {
	FileBasedStore
	MetadataStore
	ContainerBigDataStore
	FlaggableStore

	Create(id string, names []string, image, layer, metadata string, options *ContainerOptions) (*Container, error)

	// SetNames updates the list of names associated with the container
	// with the specified ID.
	SetNames(id string, names []string) error

	// Get retrieves information about a container given an ID or name.
	Get(id string) (*Container, error)

	// Exists checks if there is a container with the given ID or name.
	Exists(id string) bool

	// Delete removes the record of the container.
	Delete(id string) error

	// Wipe removes records of all containers.
	Wipe() error

	// Lookup attempts to translate a name to an ID.  Most methods do this
	// implicitly.
	Lookup(name string) (string, error)

	// Containers returns a slice enumerating the known containers.
	Containers() ([]Container, error)
}
