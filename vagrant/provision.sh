#!/usr/bin/env bash
set -xe

source /etc/os-release

case "${ID_LIKE:-${ID:-unknown}}" in
  debian)
    export DEBIAN_FRONTEND=noninteractive
    apt-get -q update
    apt-get -q -y install linux-headers-`uname -r`
    echo deb http://httpredir.debian.org/debian testing main    >  /etc/apt/sources.list
    echo deb http://httpredir.debian.org/debian testing contrib >> /etc/apt/sources.list
    apt-get -q update
    apt-get -q -y install systemd curl
    apt-get -q -y install apt make git btrfs-progs libdevmapper-dev
    apt-get -q -y install zfs-dkms zfsutils-linux
    apt-get -q -y install golang gccgo
    apt-get -q -y install bats
    ;;
  unknown)
    echo Unknown box OS, unsure of how to install required packages.
    exit 1
    ;;
esac

mkdir -p /go/src/github.com/gepis
rm -f /go/src/github.com/gepis/strge
ln -s /vagrant /go/src/github.com/gepis/strge

export GOPATH=/go
export PATH=/go/bin:${PATH}

make -C /vagrant install.tools
exit 0
