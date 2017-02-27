#!/usr/bin/env bash

# This script builds and compresses
# pilad in order to prepare a release.

set -ex

if [ -z "$1" ]; then
	echo "No release number provided!"
	exit 1
fi

make gox

oss=( linux darwin )
archs=( amd64 )
for os in "${oss[@]}"
do
	for arch in "${archs[@]}"
	do
		cd dist/$os/$arch
		tar -cvzf piladb$1.$os-$arch.tar.gz pilad
		zip -r piladb$1.$os-$arch.zip pilad
		cd ../../..
	done
done

