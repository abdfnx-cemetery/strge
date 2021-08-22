package storage

type Strge struct {
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
