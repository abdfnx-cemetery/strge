package config

import (
	"fmt"
	"os"
)

type ThinpoolOptionsConfig struct {
	AutoExtendPercent string `toml:"autoextend_percent"`
	// AutoExtendThreshold determines the pool extension threshold in terms
	// of percentage of pool size. For example, if threshold is 60, that
	// means when pool is 60% full, threshold has been hit.
	AutoExtendThreshold string `toml:"autoextend_threshold"`
	BaseSize string `toml:"basesize"`
	BlockSize string `toml:"blocksize"`
	DirectLvmDevice string `toml:"directlvm_device"`
	// DirectLvmDeviceForcewipes device even if device already has a
	// filesystem
	DirectLvmDeviceForce string `toml:"directlvm_device_force"`
	Fs string `toml:"fs"`
	LogLevel string `toml:"log_level"`
	// MetadataSize specifies the size of the metadata for the thinpool
	// It will be used with the `pvcreate --metadata` option.
	MetadataSize string `toml:"metadatasize"`
	MinFreeSpace string `toml:"min_free_space"`
	MkfsArg string `toml:"mkfsarg"`
	MountOpt string `toml:"mountopt"`
	Size string `toml:"size"`
	UseDeferredDeletion string `toml:"use_deferred_deletion"`
	UseDeferredRemoval string `toml:"use_deferred_removal"`
	XfsNoSpaceMaxRetries string `toml:"xfs_nospace_max_retries"`
}

type AufsOptionsConfig struct {
	// MountOpt specifies extra mount options used when mounting
	MountOpt string `toml:"mountopt"`
}

type BtrfsOptionsConfig struct {
	// MinSpace is the minimal spaces allocated to the device
	MinSpace string `toml:"min_space"`
	// Size
	Size string `toml:"size"`
}

