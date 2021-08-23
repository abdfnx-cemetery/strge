#define _GNU_SOURCE
#include <libdevmapper.h>
#include <linux/fs.h>

// FIXME: Can't we find a way to do the logging in pure Go?
extern void StorageDevmapperLogCallback(int level, char *file, int line, int dm_errno_or_class, char *str);

static void	log_cb(int level, const char *file, int line, int dm_errno_or_class, const char *f, ...) {
	char *buffer = NULL;
	va_list ap;
	int ret;

	va_start(ap, f);
	ret = vasprintf(&buffer, f, ap);
	va_end(ap);

	if (ret < 0) {
		// memory allocation failed -- should never happen?
		return;
	}

	StorageDevmapperLogCallback(level, (char *)file, line, dm_errno_or_class, buffer);
	free(buffer);
}

static void	log_with_errno_init() {
	dm_log_with_errno_init(log_cb);
}
