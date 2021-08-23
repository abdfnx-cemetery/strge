#include <stdlib.h>
#include <dirent.h>
#include <btrfs/ioctl.h>
#include <btrfs/ctree.h>

static void set_name_btrfs_ioctl_vol_args_v2(struct btrfs_ioctl_vol_args_v2* btrfs_struct, const char* value) {
    snprintf(btrfs_struct->name, BTRFS_SUBVOL_NAME_MAX, "%s", value);
}
