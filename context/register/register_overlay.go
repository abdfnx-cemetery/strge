// +build !exclude_graphdriver_overlay,linux,cgo

package register

import (
	// register the overlay graphdriver
	_ "github.com/gepis/strge/context/overlay"
)
