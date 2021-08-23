package vfs

import "github.com/gepis/strge/context/copy"

func dirCopy(srcDir, dstDir string) error {
	return copy.DirCopy(srcDir, dstDir, copy.Content, true)
}
