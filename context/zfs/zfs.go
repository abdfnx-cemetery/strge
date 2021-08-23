// +build linux freebsd

package zfs

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gepis/strge/context"
	"github.com/gepis/strge/pkg/directory"
	"github.com/gepis/strge/pkg/idtools"
	"github.com/gepis/strge/pkg/mount"
	"github.com/gepis/strge/pkg/parsers"
	"github.com/mistifyio/go-zfs"
	"github.com/opencontainers/selinux/go-selinux/label"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type zfsOptions struct {
	fsName       string
	mountPath    string
	mountOptions string
}

const defaultPerms = os.FileMode(0555)

func init() {
	context.Register("zfs", Init)
}

// Logger returns a zfs logger implementation.
type Logger struct{}

// Log wraps log message from ZFS driver with a prefix '[zfs]'.
func (*Logger) Log(cmd []string) {
	logrus.WithField("storage-driver", "zfs").Debugf("%s", strings.Join(cmd, " "))
}

// Init returns a new ZFS driver.
// It takes base mount path and an array of options which are represented as key value pairs.
// Each option is in the for key=value. 'zfs.fsname' is expected to be a valid key in the options.
func Init(base string, opt context.Options) (context.Driver, error) {
	var err error

	logger := logrus.WithField("storage-driver", "zfs")

	if _, err := exec.LookPath("zfs"); err != nil {
		logger.Debugf("zfs command is not available: %v", err)
		return nil, errors.Wrap(context.ErrPrerequisites, "the 'zfs' command is not available")
	}

	file, err := os.OpenFile("/dev/zfs", os.O_RDWR, 0600)
	if err != nil {
		logger.Debugf("cannot open /dev/zfs: %v", err)
		return nil, errors.Wrapf(context.ErrPrerequisites, "could not open /dev/zfs: %v", err)
	}
	defer file.Close()

	options, err := parseOptions(opt.DriverOptions)
	if err != nil {
		return nil, err
	}
	options.mountPath = base

	rootdir := path.Dir(base)

	if options.fsName == "" {
		err = checkRootdirFs(rootdir)
		if err != nil {
			return nil, err
		}
	}

	if options.fsName == "" {
		options.fsName, err = lookupZfsDataset(rootdir)
		if err != nil {
			return nil, err
		}
	}

	zfs.SetLogger(new(Logger))

	filesystems, err := zfs.Filesystems(options.fsName)
	if err != nil {
		return nil, fmt.Errorf("Cannot find root filesystem %s: %v", options.fsName, err)
	}

	filesystemsCache := make(map[string]bool, len(filesystems))
	var rootDataset *zfs.Dataset
	for _, fs := range filesystems {
		if fs.Name == options.fsName {
			rootDataset = fs
		}
		filesystemsCache[fs.Name] = true
	}

	if rootDataset == nil {
		return nil, fmt.Errorf("BUG: zfs get all -t filesystem -rHp '%s' should contain '%s'", options.fsName, options.fsName)
	}

	rootUID, rootGID, err := idtools.GetRootUIDGID(opt.UIDMaps, opt.GIDMaps)
	if err != nil {
		return nil, fmt.Errorf("Failed to get root uid/gid: %v", err)
	}
	if err := idtools.MkdirAllAs(base, 0700, rootUID, rootGID); err != nil {
		return nil, fmt.Errorf("Failed to create '%s': %v", base, err)
	}

	d := &Driver{
		dataset:          rootDataset,
		options:          options,
		filesystemsCache: filesystemsCache,
		uidMaps:          opt.UIDMaps,
		gidMaps:          opt.GIDMaps,
		ctr:              context.NewRefCounter(context.NewDefaultChecker()),
	}
	return context.NewNaiveDiffDriver(d, context.NewNaiveLayerIDMapUpdater(d)), nil
}

func parseOptions(opt []string) (zfsOptions, error) {
	var options zfsOptions
	options.fsName = ""
	for _, option := range opt {
		key, val, err := parsers.ParseKeyValueOpt(option)
		if err != nil {
			return options, err
		}
		key = strings.ToLower(key)
		switch key {
		case "zfs.fsname":
			options.fsName = val
		case "zfs.mountopt":
			options.mountOptions = val
		default:
			return options, fmt.Errorf("Unknown option %s", key)
		}
	}
	return options, nil
}

func lookupZfsDataset(rootdir string) (string, error) {
	var stat unix.Stat_t
	if err := unix.Stat(rootdir, &stat); err != nil {
		return "", fmt.Errorf("Failed to access '%s': %s", rootdir, err)
	}
	wantedDev := stat.Dev

	mounts, err := mount.GetMounts()
	if err != nil {
		return "", err
	}
	for _, m := range mounts {
		if err := unix.Stat(m.Mountpoint, &stat); err != nil {
			logrus.WithField("storage-driver", "zfs").Debugf("failed to stat '%s' while scanning for zfs mount: %v", m.Mountpoint, err)
			continue // may fail on fuse file systems
		}

		if stat.Dev == wantedDev && m.FSType == "zfs" {
			return m.Source, nil
		}
	}

	return "", fmt.Errorf("Failed to find zfs dataset mounted on '%s' in /proc/mounts", rootdir)
}

// Driver holds information about the driver, such as zfs dataset, options and cache.
type Driver struct {
	dataset          *zfs.Dataset
	options          zfsOptions
	sync.Mutex       // protects filesystem cache against concurrent access
	filesystemsCache map[string]bool
	uidMaps          []idtools.IDMap
	gidMaps          []idtools.IDMap
	ctr              *context.RefCounter
}

func (d *Driver) String() string {
	return "zfs"
}

// Cleanup is called on when program exits, it is a no-op for ZFS.
func (d *Driver) Cleanup() error {
	return nil
}

// Status returns information about the ZFS filesystem. It returns a two dimensional array of information
// such as pool name, dataset name, disk usage, parent quota and compression used.
// Currently it return 'Zpool', 'Zpool Health', 'Parent Dataset', 'Space Used By Parent',
// 'Space Available', 'Parent Quota' and 'Compression'.
func (d *Driver) Status() [][2]string {
	parts := strings.Split(d.dataset.Name, "/")
	pool, err := zfs.GetZpool(parts[0])

	var poolName, poolHealth string
	if err == nil {
		poolName = pool.Name
		poolHealth = pool.Health
	} else {
		poolName = fmt.Sprintf("error while getting pool information %v", err)
		poolHealth = "not available"
	}

	quota := "no"
	if d.dataset.Quota != 0 {
		quota = strconv.FormatUint(d.dataset.Quota, 10)
	}

	return [][2]string{
		{"Zpool", poolName},
		{"Zpool Health", poolHealth},
		{"Parent Dataset", d.dataset.Name},
		{"Space Used By Parent", strconv.FormatUint(d.dataset.Used, 10)},
		{"Space Available", strconv.FormatUint(d.dataset.Avail, 10)},
		{"Parent Quota", quota},
		{"Compression", d.dataset.Compression},
	}
}

// Metadata returns image/container metadata related to graph driver
func (d *Driver) Metadata(id string) (map[string]string, error) {
	return map[string]string{
		"Mountpoint": d.mountPath(id),
		"Dataset":    d.zfsPath(id),
	}, nil
}

func (d *Driver) cloneFilesystem(name, parentName string) error {
	snapshotName := fmt.Sprintf("%d", time.Now().Nanosecond())
	parentDataset := zfs.Dataset{Name: parentName}
	snapshot, err := parentDataset.Snapshot(snapshotName /*recursive */, false)
	if err != nil {
		return err
	}

	_, err = snapshot.Clone(name, map[string]string{"mountpoint": "legacy"})
	if err == nil {
		d.Lock()
		d.filesystemsCache[name] = true
		d.Unlock()
	}

	if err != nil {
		snapshot.Destroy(zfs.DestroyDeferDeletion)
		return err
	}
	return snapshot.Destroy(zfs.DestroyDeferDeletion)
}

func (d *Driver) zfsPath(id string) string {
	return d.options.fsName + "/" + id
}

func (d *Driver) mountPath(id string) string {
	return path.Join(d.options.mountPath, "graph", getMountpoint(id))
}

// CreateFromTemplate creates a layer with the same contents and parent as another layer.
func (d *Driver) CreateFromTemplate(id, template string, templateIDMappings *idtools.IDMappings, parent string, parentIDMappings *idtools.IDMappings, opts *context.CreateOpts, readWrite bool) error {
	return d.Create(id, template, opts)
}

// CreateReadWrite creates a layer that is writable for use as a container
// file system.
func (d *Driver) CreateReadWrite(id, parent string, opts *context.CreateOpts) error {
	return d.Create(id, parent, opts)
}

// Create prepares the dataset and filesystem for the ZFS driver for the given id under the parent.
func (d *Driver) Create(id, parent string, opts *context.CreateOpts) error {
	err := d.create(id, parent, opts)
	if err == nil {
		return nil
	}
	if zfsError, ok := err.(*zfs.Error); ok {
		if !strings.HasSuffix(zfsError.Stderr, "dataset already exists\n") {
			return err
		}
		// aborted build -> cleanup
	} else {
		return err
	}

	dataset := zfs.Dataset{Name: d.zfsPath(id)}
	if err := dataset.Destroy(zfs.DestroyRecursiveClones); err != nil {
		return err
	}

	// retry
	return d.create(id, parent, opts)
}

func (d *Driver) create(id, parent string, opts *context.CreateOpts) error {
	var storageOpt map[string]string
	if opts != nil {
		storageOpt = opts.StorageOpt
	}

	name := d.zfsPath(id)
	mountpoint := d.mountPath(id)
	quota, err := parseStorageOpt(storageOpt)
	if err != nil {
		return err
	}
	if parent == "" {
		var rootUID, rootGID int
		var mountLabel string
		if opts != nil {
			rootUID, rootGID, err = idtools.GetRootUIDGID(opts.UIDs(), opts.GIDs())
			if err != nil {
				return fmt.Errorf("Failed to get root uid/gid: %v", err)
			}
			mountLabel = opts.MountLabel
		}
		mountoptions := map[string]string{"mountpoint": "legacy"}
		fs, err := zfs.CreateFilesystem(name, mountoptions)
		if err == nil {
			err = setQuota(name, quota)
			if err == nil {
				d.Lock()
				d.filesystemsCache[fs.Name] = true
				d.Unlock()
			}

			if err := idtools.MkdirAllAs(mountpoint, defaultPerms, rootUID, rootGID); err != nil {
				return err
			}
			defer func() {
				if err := unix.Rmdir(mountpoint); err != nil && !os.IsNotExist(err) {
					logrus.Debugf("Failed to remove %s mount point %s: %v", id, mountpoint, err)
				}
			}()

			mountOpts := label.FormatMountLabel(d.options.mountOptions, mountLabel)

			if err := mount.Mount(name, mountpoint, "zfs", mountOpts); err != nil {
				return errors.Wrap(err, "error creating zfs mount")
			}
			defer func() {
				if err := unix.Unmount(mountpoint, unix.MNT_DETACH); err != nil {
					logrus.Warnf("Failed to unmount %s mount %s: %v", id, mountpoint, err)
				}
			}()

			if err := os.Chmod(mountpoint, defaultPerms); err != nil {
				return errors.Wrap(err, "error setting permissions on zfs mount")
			}

			// this is our first mount after creation of the filesystem, and the root dir may still have root
			// permissions instead of the remapped root uid:gid (if user namespaces are enabled):
			if err := os.Chown(mountpoint, rootUID, rootGID); err != nil {
				return errors.Wrapf(err, "modifying zfs mountpoint (%s) ownership", mountpoint)
			}

		}
		return err
	}
	err = d.cloneFilesystem(name, d.zfsPath(parent))
	if err == nil {
		err = setQuota(name, quota)
	}
	return err
}

func parseStorageOpt(storageOpt map[string]string) (string, error) {
	// Read size to change the disk quota per container
	for k, v := range storageOpt {
		key := strings.ToLower(k)
		switch key {
		case "size":
			return v, nil
		default:
			return "0", fmt.Errorf("Unknown option %s", key)
		}
	}
	return "0", nil
}

func setQuota(name string, quota string) error {
	if quota == "0" {
		return nil
	}
	fs, err := zfs.GetDataset(name)
	if err != nil {
		return err
	}
	return fs.SetProperty("quota", quota)
}

// Remove deletes the dataset, filesystem and the cache for the given id.
func (d *Driver) Remove(id string) error {
	name := d.zfsPath(id)
	dataset := zfs.Dataset{Name: name}
	err := dataset.Destroy(zfs.DestroyRecursive)
	if err == nil {
		d.Lock()
		delete(d.filesystemsCache, name)
		d.Unlock()
	}
	return err
}

// Get returns the mountpoint for the given id after creating the target directories if necessary.
func (d *Driver) Get(id string, options context.MountOpts) (_ string, retErr error) {

	mountpoint := d.mountPath(id)
	if count := d.ctr.Increment(mountpoint); count > 1 {
		return mountpoint, nil
	}
	defer func() {
		if retErr != nil {
			if c := d.ctr.Decrement(mountpoint); c <= 0 {
				if mntErr := unix.Unmount(mountpoint, 0); mntErr != nil {
					logrus.WithField("storage-driver", "zfs").Errorf("Error unmounting %v: %v", mountpoint, mntErr)
				}
				if rmErr := unix.Rmdir(mountpoint); rmErr != nil && !os.IsNotExist(rmErr) {
					logrus.WithField("storage-driver", "zfs").Debugf("Failed to remove %s: %v", id, rmErr)
				}

			}
		}
	}()

	// In the case of a read-only mount we first mount read-write so we can set the
	// correct permissions on the mount point and remount read-only afterwards.
	remountReadOnly := false
	mountOptions := d.options.mountOptions
	if len(options.Options) > 0 {
		var newOptions []string
		for _, option := range options.Options {
			if option == "ro" {
				// Filter out read-only mount option but remember for later remounting.
				remountReadOnly = true
			} else {
				newOptions = append(newOptions, option)
			}
		}
		mountOptions = strings.Join(newOptions, ",")
	}

	filesystem := d.zfsPath(id)
	opts := label.FormatMountLabel(mountOptions, options.MountLabel)
	logrus.WithField("storage-driver", "zfs").Debugf(`mount("%s", "%s", "%s")`, filesystem, mountpoint, opts)

	rootUID, rootGID, err := idtools.GetRootUIDGID(d.uidMaps, d.gidMaps)
	if err != nil {
		return "", err
	}
	// Create the target directories if they don't exist
	if err := idtools.MkdirAllAs(mountpoint, 0755, rootUID, rootGID); err != nil {
		return "", err
	}

	if err := mount.Mount(filesystem, mountpoint, "zfs", opts); err != nil {
		return "", errors.Wrap(err, "error creating zfs mount")
	}

	if remountReadOnly {
		opts = label.FormatMountLabel("remount,ro", options.MountLabel)
		if err := mount.Mount(filesystem, mountpoint, "zfs", opts); err != nil {
			return "", errors.Wrap(err, "error remounting zfs mount read-only")
		}
	}

	return mountpoint, nil
}

// Put removes the existing mountpoint for the given id if it exists.
func (d *Driver) Put(id string) error {
	mountpoint := d.mountPath(id)
	if count := d.ctr.Decrement(mountpoint); count > 0 {
		return nil
	}

	logger := logrus.WithField("storage-driver", "zfs")

	logger.Debugf(`unmount("%s")`, mountpoint)

	if err := unix.Unmount(mountpoint, unix.MNT_DETACH); err != nil {
		logger.Warnf("Failed to unmount %s mount %s: %v", id, mountpoint, err)
	}
	if err := unix.Rmdir(mountpoint); err != nil && !os.IsNotExist(err) {
		logger.Debugf("Failed to remove %s mount point %s: %v", id, mountpoint, err)
	}

	return nil
}

// ReadWriteDiskUsage returns the disk usage of the writable directory for the ID.
// For ZFS, it queries the full mount path for this ID.
func (d *Driver) ReadWriteDiskUsage(id string) (*directory.DiskUsage, error) {
	return directory.Usage(d.mountPath(id))
}

// Exists checks to see if the cache entry exists for the given id.
func (d *Driver) Exists(id string) bool {
	d.Lock()
	defer d.Unlock()
	return d.filesystemsCache[d.zfsPath(id)]
}

// AdditionalImageStores returns additional image stores supported by the driver
func (d *Driver) AdditionalImageStores() []string {
	return nil
}