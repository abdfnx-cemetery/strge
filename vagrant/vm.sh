#!/usr/bin/env bash

set -e
export PKG="github.com/gepis/strge"
MACHINE="debian"

if test -z "$VAGRANT_PROVIDER" ; then
	if lsmod | grep -q '^vboxdrv ' ; then
		VAGRANT_PROVIDER=virtualbox
	elif lsmod | grep -q '^kvm ' ; then
		VAGRANT_PROVIDER=libvirt
	fi
fi

export VAGRANT_PROVIDER=${VAGRANT_PROVIDER:-libvirt}

if ${IN_VAGRANT_MACHINE:-false} ; then
	unset AUTO_GOPATH
	export GOPATH=/go
	export PATH=${GOPATH}/bin:/go/src/${PKG}/vendor/src/github.com/golang/lint/golint:${PATH}
	sudo modprobe aufs || true
	sudo modprobe zfs || true
	"$@"
else
	vagrant up --provider ${VAGRANT_PROVIDER}
	vagrant reload ${MACHINE}
	vagrant ssh ${MACHINE} -c "cd /go/src/${PKG}; IN_VAGRANT_MACHINE=true sudo -E $0 $*"
	vagrant ssh ${MACHINE} -c "sudo poweroff &"
fi
