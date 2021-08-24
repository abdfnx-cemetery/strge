#!/usr/bin/env bash
#
# Work around, based on one described in https://github.com/golang/go/issues/15628
#
addflags=
for arg in "$@" ; do
	if test -d "$arg"/github.com/gepis/strge/vendor ; then
		addflags="$addflags -I $arg/github.com/gepis/strge/vendor"
	fi
done
exec gccgo $addflags "$@"
