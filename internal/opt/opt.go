package opt

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

var (
	alphaRegexp  = regexp.MustCompile(`[a-zA-Z]`)
	domainRegexp = regexp.MustCompile(`^(:?(:?[a-zA-Z0-9]|(:?[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9]))(:?\.(:?[a-zA-Z0-9]|(:?[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])))*)\.?\s*$`)
)

// ListOpt holds a list of values and a validation function.
type ListOpt struct {
	values    *[]string
	validator ValidatorFctType
}

// NewListOpt creates a new ListOpt with the specified validator.
func NewListOpt(validator ValidatorFctType) ListOpt {
	var values []string
	return *NewListOptRef(&values, validator)
}

// NewListOptRef creates a new ListOpt with the specified values and validator.
func NewListOptRef(values *[]string, validator ValidatorFctType) *ListOpt {
	return &ListOpt{
		values:    values,
		validator: validator,
	}
}

func (opt *ListOpt) String() string {
	return fmt.Sprintf("%v", []string((*opt.values)))
}

// Set validates if needed the input value and adds it to the
// internal slice.
func (opt *ListOpt) Set(value string) error {
	if opt.validator != nil {
		v, err := opt.validator(value)
		if err != nil {
			return err
		}

		value = v
	}

	(*opt.values) = append((*opt.values), value)
	return nil
}

// Delete removes the specified element from the slice.
func (opt *ListOpt) Delete(key string) {
	for i, k := range *opt.values {
		if k == key {
			(*opt.values) = append((*opt.values)[:i], (*opt.values)[i+1:]...)
			return
		}
	}
}

// GetMap returns the content of values in a map in order to avoid
// duplicates.
func (opt *ListOpt) GetMap() map[string]struct{} {
	ret := make(map[string]struct{})
	for _, k := range *opt.values {
		ret[k] = struct{}{}
	}

	return ret
}

// GetAll returns the values of slice.
func (opt *ListOpt) GetAll() []string {
	return (*opt.values)
}

// GetAllOrEmpty returns the values of the slice
// or an empty slice when there are no values.
func (opt *ListOpt) GetAllOrEmpty() []string {
	v := *opt.values
	if v == nil {
		return make([]string, 0)
	}

	return v
}

// Get checks the existence of the specified key.
func (opt *ListOpt) Get(key string) bool {
	for _, k := range *opt.values {
		if k == key {
			return true
		}
	}

	return false
}

// Len returns the amount of element in the slice.
func (opt *ListOpt) Len() int {
	return len((*opt.values))
}

// Type returns a string name for this Option type
func (opt *ListOpt) Type() string {
	return "list"
}

// NamedOption is an interface that list and map options
// with names implement.
type NamedOption interface {
	Name() string
}

// NamedListOpt is a ListOpt with a configuration name.
// This struct is useful to keep reference to the assigned
// field name in the internal configuration struct.
type NamedListOpt struct {
	name string
	ListOpt
}

var _ NamedOption = &NamedListOpt{}

// NewNamedListOptRef creates a reference to a new NamedListOpt struct.
func NewNamedListOptRef(name string, values *[]string, validator ValidatorFctType) *NamedListOpt {
	return &NamedListOpt{
		name:     name,
		ListOpt: *NewListOptRef(values, validator),
	}
}

// Name returns the name of the NamedListOpt in the configuration.
func (o *NamedListOpt) Name() string {
	return o.name
}

//MapOpt holds a map of values and a validation function.
type MapOpt struct {
	values    map[string]string
	validator ValidatorFctType
}

// Set validates if needed the input value and add it to the
// internal map, by splitting on '='.
func (opt *MapOpt) Set(value string) error {
	if opt.validator != nil {
		v, err := opt.validator(value)
		if err != nil {
			return err
		}

		value = v
	}

	vals := strings.SplitN(value, "=", 2)
	if len(vals) == 1 {
		(opt.values)[vals[0]] = ""
	} else {
		(opt.values)[vals[0]] = vals[1]
	}

	return nil
}

// GetAll returns the values of MapOpt as a map.
func (opt *MapOpt) GetAll() map[string]string {
	return opt.values
}

func (opt *MapOpt) String() string {
	return fmt.Sprintf("%v", map[string]string((opt.values)))
}

// Type returns a string name for this Option type
func (opt *MapOpt) Type() string {
	return "map"
}

// NewMapOpt creates a new MapOpt with the specified map of values and a validator.
func NewMapOpt(values map[string]string, validator ValidatorFctType) *MapOpt {
	if values == nil {
		values = make(map[string]string)
	}

	return &MapOpt{
		values:    values,
		validator: validator,
	}
}

// NamedMapOpt is a MapOpt struct with a configuration name.
// This struct is useful to keep reference to the assigned
// field name in the internal configuration struct.
type NamedMapOpt struct {
	name string
	MapOpt
}

var _ NamedOption = &NamedMapOpt{}

// NewNamedMapOpt creates a reference to a new NamedMapOpt struct.
func NewNamedMapOpt(name string, values map[string]string, validator ValidatorFctType) *NamedMapOpt {
	return &NamedMapOpt{
		name:    name,
		MapOpt: *NewMapOpt(values, validator),
	}
}

// Name returns the name of the NamedMapOpt in the configuration.
func (o *NamedMapOpt) Name() string {
	return o.name
}

// ValidatorFctType defines a validator function that returns a validated string and/or an error.
type ValidatorFctType func(val string) (string, error)

// ValidatorFctListType defines a validator function that returns a validated list of string and/or an error
type ValidatorFctListType func(val string) ([]string, error)

// ValidateIPAddress validates an Ip address.
func ValidateIPAddress(val string) (string, error) {
	var ip = net.ParseIP(strings.TrimSpace(val))
	if ip != nil {
		return ip.String(), nil
	}

	return "", fmt.Errorf("%s is not an ip address", val)
}

// ValidateDNSSearch validates domain for resolvconf search configuration.
// A zero length domain is represented by a dot (.).
func ValidateDNSSearch(val string) (string, error) {
	if val = strings.Trim(val, " "); val == "." {
		return val, nil
	}

	return validateDomain(val)
}

func validateDomain(val string) (string, error) {
	if alphaRegexp.FindString(val) == "" {
		return "", fmt.Errorf("%s is not a valid domain", val)
	}

	ns := domainRegexp.FindSubmatch([]byte(val))
	if len(ns) > 0 && len(ns[1]) < 255 {
		return string(ns[1]), nil
	}

	return "", fmt.Errorf("%s is not a valid domain", val)
}

// ValidateLabel validates that the specified string is a valid label, and returns it.
// Labels are in the form on key=value.
func ValidateLabel(val string) (string, error) {
	if strings.Count(val, "=") < 1 {
		return "", fmt.Errorf("bad attribute format: %s", val)
	}
	return val, nil
}

// ValidateSysctl validates a sysctl and returns it.
func ValidateSysctl(val string) (string, error) {
	validSysctlMap := map[string]bool{
		"kernel.msgmax":          true,
		"kernel.msgmnb":          true,
		"kernel.msgmni":          true,
		"kernel.sem":             true,
		"kernel.shmall":          true,
		"kernel.shmmax":          true,
		"kernel.shmmni":          true,
		"kernel.shm_rmid_forced": true,
	}

	validSysctlPrefixes := []string{
		"net.",
		"fs.mqueue.",
	}

	arr := strings.Split(val, "=")
	if len(arr) < 2 {
		return "", fmt.Errorf("sysctl '%s' is not allowed", val)
	}

	if validSysctlMap[arr[0]] {
		return val, nil
	}

	for _, vp := range validSysctlPrefixes {
		if strings.HasPrefix(arr[0], vp) {
			return val, nil
		}
	}

	return "", fmt.Errorf("sysctl '%s' is not allowed", val)
}

// FilterOpt is a flag type for validating filters
type FilterOpt struct {
	filter Args
}

// NewFilterOpt returns a new FilterOpt
func NewFilterOpt() FilterOpt {
	return FilterOpt{filter: NewArgs()}
}

func (o *FilterOpt) String() string {
	repr, err := ToParam(o.filter)
	if err != nil {
		return "invalid filters"
	}

	return repr
}

// Set sets the value of the opt by parsing the command line value
func (o *FilterOpt) Set(value string) error {
	var err error
	o.filter, err = ParseFlag(value, o.filter)
	return err
}

// Type returns the option type
func (o *FilterOpt) Type() string {
	return "filter"
}

// Value returns the value of this option
func (o *FilterOpt) Value() Args {
	return o.filter
}
