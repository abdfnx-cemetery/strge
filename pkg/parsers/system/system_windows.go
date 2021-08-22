package system

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// GetSystem gets the name of the current operating system.
func GetSystem() (string, error) {
	var h windows.Handle

	// Default return value
	ret := "Unknown Operating System"

	if err := windows.RegOpenKeyEx(windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\`),
		0,
		windows.KEY_READ,
		&h); err != nil {
		return ret, err
	}

	defer windows.RegCloseKey(h)

	var buf [1 << 10]uint16
	var typ uint32
	n := uint32(len(buf) * 2) // api expects array of bytes, not uint16

	if err := windows.RegQueryValueEx(h,
		windows.StringToUTF16Ptr("ProductName"),
		nil,
		&typ,
		(*byte)(unsafe.Pointer(&buf[0])),
		&n); err != nil {
		return ret, err
	}

	ret = windows.UTF16ToString(buf[:])

	return ret, nil
}

// IsContainerized returns true if we are running inside a container.
// No-op on Windows, always returns false.
func IsContainerized() (bool, error) {
	return false, nil
}
