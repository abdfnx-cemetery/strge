// +build !windows

package constants

// LCOWSupported returns true if Linux containers on Windows are supported.
func LCOWSupported() bool {
	return false
}
