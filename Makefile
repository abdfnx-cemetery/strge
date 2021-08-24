export GO111MODULE=off
export GOPROXY=https://proxy.golang.org

.PHONY: \
	all \
	binary \
	clean \
	cross \
	default \
	gccgo \
	help \
	install.tools \
	local-binary \
	local-cross \
	local-gccgo \
	local-validate \
	validate \
	vendor

PACKAGE := github.com/gepis/strge
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_BRANCH_CLEAN := $(shell echo $(GIT_BRANCH) | sed -e "s/[^[:alnum:]]/-/g")
NATIVETAGS :=
AUTOTAGS := $(shell ./scripts/btrfs_tag.sh) $(shell ./scripts/libdm_tag.sh) $(shell ./scripts/libsubid_tag.sh)
BUILDFLAGS := -tags "$(AUTOTAGS) $(TAGS)" $(FLAGS)
GO ?= go

# Go module support: set `-mod=vendor` to use the vendored sources
ifeq ($(shell $(GO) help mod >/dev/null 2>&1 && echo true), true)
	GO:=GO111MODULE=on $(GO)
	MOD_VENDOR=-mod=vendor
endif

RUNINVM := vagrant/vm.sh

default all: local-binary local-validate local-cross local-gccgo

clean: ## remove all built files
	$(RM) -f gepis-strge gepis-strge.*

sources := $(wildcard *.go core/gepis-strge/*.go drivers/*.go drivers/*/*.go pkg/*/*.go pkg/*/*/*.go)
gepis-strge: $(sources) ## build using gc on the host
	$(GO) build $(MOD_VENDOR) -compiler gc $(BUILDFLAGS) ./core

binary local-binary: gepis-strge

local-gccgo: ## build using gccgo on the host
	GCCGO=$(PWD)/scripts/gccgo-wrapper.sh $(GO) build $(MOD_VENDOR) -compiler gccgo $(BUILDFLAGS) -o gepis-strge.gccgo ./core

local-cross: ## cross build the binaries for arm, darwin, and\nfreebsd
	@for target in linux/amd64 linux/386 linux/arm linux/arm64 linux/ppc64 linux/ppc64le darwin/amd64 windows/amd64 ; do \
		os=`echo $${target} | cut -f1 -d/` ; \
		arch=`echo $${target} | cut -f2 -d/` ; \
		suffix=$${os}.$${arch} ; \
		echo env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} $(GO) build $(MOD_VENDOR) -compiler gc -tags \"$(NATIVETAGS) $(TAGS)\" $(FLAGS) -o gepis-strge.$${suffix} ./core ; \
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} $(GO) build $(MOD_VENDOR) -compiler gc -tags "$(NATIVETAGS) $(TAGS)" $(FLAGS) -o gepis-strge.$${suffix} ./core || exit 1 ; \
	done

cross: ## cross build the binaries for arm, darwin, and\nfreebsd using VMs
	$(RUNINVM) make local-$@

gccgo: ## build using gccgo using VMs
	$(RUNINVM) make local-$@

validate: ## validate DCO, gofmt, ./pkg/ isolation, golint,\ngo vet and vendor using VMs
	$(RUNINVM) make local-$@

help: ## this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-z A-Z_-]+:.*?## / {gsub(" ",",",$$1);gsub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-21s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

vendor-in-container:
	gepis run --privileged --rm --env HOME=/root -v `pwd`:/src -w /src golang make vendor

vendor:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify
