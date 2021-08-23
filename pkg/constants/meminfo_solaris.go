// +build solaris,cgo

package constants

import (
	"fmt"
	"unsafe"
)

import "C"

// Get the system memory info using sysconf same as prtconf
func getTotalMem() int64 {
	pagesize := C.sysconf(C._SC_PAGESIZE)
	npages := C.sysconf(C._SC_PHYS_PAGES)
	return int64(pagesize * npages)
}

func getFreeMem() int64 {
	pagesize := C.sysconf(C._SC_PAGESIZE)
	npages := C.sysconf(C._SC_AVPHYS_PAGES)
	return int64(pagesize * npages)
}

// ReadMemInfo retrieves memory statistics of the host system and returns a
//  MemInfo type.
func ReadMemInfo() (*MemInfo, error) {

	ppKernel := C.getPpKernel()
	MemTotal := getTotalMem()
	MemFree := getFreeMem()
	SwapTotal, SwapFree, err := getSysSwap()

	if ppKernel < 0 || MemTotal < 0 || MemFree < 0 || SwapTotal < 0 ||
		SwapFree < 0 {
		return nil, fmt.Errorf("error getting system memory info %v\n", err)
	}

	meminfo := &MemInfo{}
	// Total memory is total physical memory less than memory locked by kernel
	meminfo.MemTotal = MemTotal - int64(ppKernel)
	meminfo.MemFree = MemFree
	meminfo.SwapTotal = SwapTotal
	meminfo.SwapFree = SwapFree

	return meminfo, nil
}

func getSysSwap() (int64, int64, error) {
	var tSwap int64
	var fSwap int64
	var diskblksPerPage int64
	num, err := C.swapctl(C.SC_GETNSWP, nil)
	if err != nil {
		return -1, -1, err
	}
	st := C.allocSwaptable(num)
	_, err = C.swapctl(C.SC_LIST, unsafe.Pointer(st))
	if err != nil {
		C.freeSwaptable(st)
		return -1, -1, err
	}

	diskblksPerPage = int64(C.sysconf(C._SC_PAGESIZE) >> C.DEV_BSHIFT)
	for i := 0; i < int(num); i++ {
		swapent := C.getSwapEnt(&st.swt_ent[0], C.int(i))
		tSwap += int64(swapent.ste_pages) * diskblksPerPage
		fSwap += int64(swapent.ste_free) * diskblksPerPage
	}
	C.freeSwaptable(st)
	return tSwap, fSwap, nil
}
