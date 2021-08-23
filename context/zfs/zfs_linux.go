package zfs

import (
	"github.com/gepis/strge/context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func checkRootdirFs(rootDir string) error {
	fsMagic, err := context.GetFSMagic(rootDir)
	if err != nil {
		return err
	}
	backingFS := "unknown"
	if fsName, ok := context.FsNames[fsMagic]; ok {
		backingFS = fsName
	}

	if fsMagic != context.FsMagicZfs {
		logrus.WithField("root", rootDir).WithField("backingFS", backingFS).WithField("storage-driver", "zfs").Error("No zfs dataset found for root")
		return errors.Wrapf(context.ErrPrerequisites, "no zfs dataset found for rootdir '%s'", rootDir)
	}

	return nil
}

func getMountpoint(id string) string {
	return id
}
