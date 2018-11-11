#!/bin/bash

# This scripts builds RPMs from the current git head.
# The script needs be run from the root of the repository
# NOTE: RPMs are built only for EL7 (CentOS7) distributions.

set -e

##
## Set up build environment
##
RESULTDIR=${RESULTDIR:-$PWD/rpms}
BUILDDIR=$PWD/$(mktemp -d nightlyrpmXXXXXX)

BASEDIR=$(dirname "$0")
GBLOCKRESTCLONE=$(realpath "$BASEDIR/..")

yum -y install make mock rpm-build golang

export GOPATH=$BUILDDIR/go
mkdir -p "$GOPATH"/{bin,pkg,src}
export PATH=$GOPATH/bin:$PATH

GBLOCKRESTSRC=$GOPATH/src/github.com/aravindavk/gluster-block-restapi
mkdir -p "$GOPATH/src/github.com/aravindavk/"
ln -s "$GBLOCKRESTCLONE" "$GBLOCKRESTSRC"

INSTALL_GOMETALINTER=no "$GBLOCKRESTSRC/scripts/install-reqs.sh"

##
## Prepare gluster-block-restapi archives and specfile for building RPMs
##
pushd "$GBLOCKRESTSRC"

VERSION=$(./scripts/pkg-version --version)
RELEASE=$(./scripts/pkg-version --release)
FULL_VERSION=$(./scripts/pkg-version --full)

# Create a vendored dist archive
DISTDIR=$BUILDDIR SIGN=no make dist-vendor

# Copy over specfile to the BUILDDIR and modify it to use the current Git HEAD versions
cp ./extras/rpms/* "$BUILDDIR"

popd #GBLOCKRESTSRC

pushd "$BUILDDIR"

DISTARCHIVE="gluster-block-restapi-$FULL_VERSION-vendor.tar.xz"
SPEC=gluster-block-restapi.spec
sed -i -E "
# Use bundled always
s/with_bundled 0/with_bundled 1/;
# Replace version with HEAD version
s/%global gluster_block_restapi_ver[[:space:]]+(.+)$/%global gluster_block_restapi_ver $VERSION/;
# Replace release with proper release
s/%global gluster_block_restapi_rel[[:space:]]+(.+)$/%global gluster_block_restapi_rel $RELEASE/;
# Replace Source0 with generated archive
s/^Source0:[[:space:]]+.*.tar.xz/Source0: $DISTARCHIVE/;
" $SPEC

# Create SRPM
mkdir -p rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
cp "$BUILDDIR/$DISTARCHIVE" rpmbuild/SOURCES
cp $SPEC rpmbuild/SPECS
SRPM=$(rpmbuild --define "_topdir $PWD/rpmbuild" -bs rpmbuild/SPECS/$SPEC | cut -d\  -f2)

# Build RPM from SRPM using mock
mkdir -p "$RESULTDIR"
/usr/bin/mock -r epel-7-x86_64 --resultdir="$RESULTDIR" --rebuild "$SRPM"

popd #BUILDDIR

## Cleanup
rm -rf "$BUILDDIR"
