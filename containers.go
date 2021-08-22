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

type containerStore struct {
	lockfile   Locker
	dir        string
	containers []*Container
	idindex    *truncindex.TruncIndex
	byid       map[string]*Container
	bylayer    map[string]*Container
	byname     map[string]*Container
	loadMut    sync.Mutex
}

func copyContainer(c *Container) *Container {
	return &Container{
		ID:             c.ID,
		Names:          copyStringSlice(c.Names),
		ImageID:        c.ImageID,
		LayerID:        c.LayerID,
		Metadata:       c.Metadata,
		BigDataNames:   copyStringSlice(c.BigDataNames),
		BigDataSizes:   copyStringInt64Map(c.BigDataSizes),
		BigDataDigests: copyStringDigestMap(c.BigDataDigests),
		Created:        c.Created,
		UIDMap:         copyIDMap(c.UIDMap),
		GIDMap:         copyIDMap(c.GIDMap),
		Flags:          copyStringInterfaceMap(c.Flags),
	}
}

func (c *Container) MountLabel() string {
	if label, ok := c.Flags["MountLabel"].(string); ok {
		return label
	}

	return ""
}

func (c *Container) ProcessLabel() string {
	if label, ok := c.Flags["ProcessLabel"].(string); ok {
		return label
	}

	return ""
}

func (c *Container) MountOpts() []string {
	switch c.Flags["MountOpts"].(type) {
		case []string:
			return c.Flags["MountOpts"].([]string)
		case []interface{}:
			var mountOpts []string
			for _, v := range c.Flags["MountOpts"].([]interface{}) {
				if flag, ok := v.(string); ok {
					mountOpts = append(mountOpts, flag)
				}
			}
			return mountOpts
		default:
			return nil
	}
}

func (r *containerStore) Containers() ([]Container, error) {
	containers := make([]Container, len(r.containers))
	
	return containers, nil
}
