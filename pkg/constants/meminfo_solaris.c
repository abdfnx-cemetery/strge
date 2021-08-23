#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lkstat
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <kstat.h>
#include <sys/swap.h>
#include <sys/param.h>
struct swaptable *allocSwaptable(int num) {
	struct swaptable *st;
	struct swapent *swapent;
	st = (struct swaptable *)malloc(num * sizeof(swapent_t) + sizeof (int));
	swapent = st->swt_ent;
	for (int i = 0; i < num; i++,swapent++) {
		swapent->ste_path = (char *)malloc(MAXPATHLEN * sizeof (char));
	}
	st->swt_n = num;
	return st;
}

void freeSwaptable (struct swaptable *st) {
	struct swapent *swapent = st->swt_ent;
	for (int i = 0; i < st->swt_n; i++,swapent++) {
		free(swapent->ste_path);
	}
	free(st);
}

swapent_t getSwapEnt(swapent_t *ent, int i) {
	return ent[i];
}

int64_t getPpKernel() {
	int64_t pp_kernel = 0;
	kstat_ctl_t *ksc;
	kstat_t *ks;
	kstat_named_t *knp;
	kid_t kid;

	if ((ksc = kstat_open()) == NULL) {
		return -1;
	}

	if ((ks = kstat_lookup(ksc, "unix", 0, "system_pages")) == NULL) {
		return -1;
	}

	if (((kid = kstat_read(ksc, ks, NULL)) == -1) ||
	    ((knp = kstat_data_lookup(ks, "pp_kernel")) == NULL)) {
		return -1;
	}

	switch (knp->data_type) {
        case KSTAT_DATA_UINT64:
            pp_kernel = knp->value.ui64;
            break;
        case KSTAT_DATA_UINT32:
            pp_kernel = knp->value.ui32;
            break;
	}

	pp_kernel *= sysconf(_SC_PAGESIZE);
	return (pp_kernel > 0 ? pp_kernel : -1);
}
