// +build linux,!btrfs_noversion,cgo

package btrfs

import "C"

func btrfsBuildVersion() string {
	return string(C.BTRFS_BUILD_VERSION)
}

func btrfsLibVersion() int {
	return int(C.BTRFS_LIB_VERSION)
}
