struct dm_ioctl {
	uint32_t version[3];
	uint32_t data_size;
	uint32_t data_start;
	uint32_t target_count;
	int32_t open_count;
	uint32_t flags;
	uint32_t event_nr;
	uint32_t padding;
	uint64_t dev;
	char name[DM_NAME_LEN];
	char uuid[DM_UUID_LEN];
	char data[7];
};

struct target {
	uint64_t start;
	uint64_t length;
	char *type;
	char *params;
	struct target *next;
};

typedef enum {
	DM_ADD_NODE_ON_RESUME, /* add /dev/mapper node with dmsetup resume */
	DM_ADD_NODE_ON_CREATE  /* add /dev/mapper node with dmsetup create */
} dm_add_node_t;

struct dm_task {
	int type;
	char *dev_name;
	char *mangled_dev_name;

	struct target *head, *tail;

	int read_only;
	uint32_t event_nr;
	int major;
	int minor;
	int allow_default_major_fallback;
	uid_t uid;
	gid_t gid;
	mode_t mode;
	uint32_t read_ahead;
	uint32_t read_ahead_flags;
	union {
		struct dm_ioctl *v4;
	} dmi;

	char *newname;
	char *message;
	char *geometry;
	uint64_t sector;
	int no_flush;
	int no_open_count;
	int skip_lockfs;
	int query_inactive_table;
	int suppress_identical_reload;
	dm_add_node_t add_node;
	uint64_t existing_table_size;
	int cookie_set;
	int new_uuid;
	int secure_data;
	int retry_remove;
	int enable_checks;
	int expected_errno;

	char *uuid;
	char *mangled_uuid;
};
