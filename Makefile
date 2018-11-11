# Setting up standard path variables similar to autoconf
# The defaults are taken based on
# https://www.gnu.org/prep/standards/html_node/Directory-Variables.html
# and
# https://fedoraproject.org/wiki/Packaging:RPMMacros?rd=Packaging/RPMMacros

PREFIX ?= /usr/local

BASE_PREFIX = $(PREFIX)
ifeq ($(PREFIX), /usr)
BASE_PREFIX = ""
endif

EXEC_PREFIX ?= $(PREFIX)

BINDIR ?= $(EXEC_PREFIX)/bin
SBINDIR ?= $(EXEC_PREFIX)/sbin

DATADIR ?= $(PREFIX)/share
LOCALSTATEDIR ?= $(BASE_PREFIX)/var/lib
LOGDIR ?= $(BASE_PREFIX)/var/log

SYSCONFDIR ?= $(BASE_PREFIX)/etc
RUNDIR ?= $(BASE_PREFIX)/var/run


PROJ_NAME = gluster-block-restapi
BINNAME = glusterblockrestd

BUILDDIR = build

RESTD_BIN = $(BINNAME)
RESTD_BUILD = $(BUILDDIR)/$(BINNAME)
RESTD_INSTALL = $(DESTDIR)$(SBINDIR)/$(BINNAME)
RESTD_SERVICE_BUILD = $(BUILDDIR)/$(BINNAME).service
RESTD_SERVICE_INSTALL = $(DESTDIR)/usr/lib/systemd/system/$(BINNAME).service
RESTD_CONF_INSTALL = $(DESTDIR)$(SYSCONFDIR)/$(PROJ_NAME)

RESTD_LOGDIR = $(LOGDIR)/$(PROJ_NAME)
RESTD_RUNDIR = $(RUNDIR)/$(PROJ_NAME)

DEPENV ?=

FASTBUILD ?= yes

.PHONY: all build binaries check check-go check-reqs install vendor-update vendor-install verify release check-protoc $(BINNAME) test dist dist-vendor gen-service gen-version

all: build

build: check-go check-reqs vendor-install $(BINNAME)
check: check-go check-reqs

check-go:
	@./scripts/check-go.sh
	@echo

check-reqs:
	@./scripts/check-reqs.sh
	@echo

$(BINNAME): gen-service
	FASTBUILD=$(FASTBUILD) BASE_PREFIX=$(BASE_PREFIX) \
		CONFFILE=${SYSCONFDIR}/gluster-block-restapi/config.toml \
		./scripts/build.sh $(BINNAME)
	@echo

install:
	install -D $(RESTD_BUILD) $(RESTD_INSTALL)
	install -D -m 0644 $(RESTD_SERVICE_BUILD) $(RESTD_SERVICE_INSTALL)
	install -D -m 0600 ./extras/conf/config.toml.sample $(RESTD_CONF_INSTALL)/config.toml
	@echo

vendor-update:
	@echo Updating vendored packages
	@$(DEPENV) dep ensure -v -update
	@echo

vendor-install:
	@echo Installing vendored packages
	@$(DEPENV) dep ensure -v -vendor-only
	@echo

test: check-reqs
	@./scripts/pre-commit.sh
	@./scripts/gometalinter-tests.sh
	@echo

dist: gen-version
	@DISTDIR=$(DISTDIR) SIGN=$(SIGN) ./scripts/dist.sh
	@rm -f ./VERSION ./GIT_SHA_FULL

dist-vendor: vendor-install gen-version
	@VENDOR=yes DISTDIR=$(DISTDIR) SIGN=$(SIGN) ./scripts/dist.sh
	@rm -f ./VERSION ./GIT_SHA_FULL

gen-service:
	SBINDIR=$(SBINDIR) SYSCONFDIR=$(SYSCONFDIR) \
		./scripts/gen-service.sh

gen-version:
	@git describe --tags --always --match "v[0-9]*" > ./VERSION
	@git rev-parse HEAD > ./GIT_SHA_FULL
