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

	for i := range r.containers {
		containers[i] = *copyContainer(r.containers[i])
	}
	
	return containers, nil
}

func (r *containerStore) containerspath() string {
	return filepath.Join(r.dir, "containers.json")
}

func (r *containerStore) datadir(id string) string {
	return filepath.Join(r.dir, id)
}

func (r *containerStore) datapath(id, key string) string {
	return filepath.Join(r.datadir(id), makeBigDataBaseName(key))
}

func (r *containerStore) Load() error {
	needSave := false
	rpath := r.containerspath()
	data, err := ioutil.ReadFile(rpath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	containers := []*Container{}
	layers := make(map[string]*Container)
	idlist := []string{}
	ids := make(map[string]*Container)
	names := make(map[string]*Container)

	if err = json.Unmarshal(data, &containers); len(data) == 0 || err == nil {
		idlist = make([]string, 0, len(containers))
		for n, container := range containers {
			idlist = append(idlist, container.ID)
			ids[container.ID] = containers[n]
			layers[container.LayerID] = containers[n]
			for _, name := range container.Names {
				if conflict, ok := names[name]; ok {
					r.removeName(conflict, name)
					needSave = true
				}

				names[name] = containers[n]
			}
		}
	}

	r.containers = containers
	r.idindex = truncindex.NewTruncIndex(idlist)
	r.byid = ids
	r.bylayer = layers
	r.byname = names

	if needSave {
		return r.Save()
	}

	return nil
}
