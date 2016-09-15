GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
LINTIGNOREDEPS='vendor/.+\.go'
TARGET_ONLY_PKGS=$(shell go list ./... 2> /dev/null | grep -v "/misc/" | grep -v "/vendor/")
INTERNAL_BIN=.bin
HAVE_GLIDE:=$(shell which glide)
HAVE_GOLINT:=$(shell which golint)
HAVE_GOCYCLO:=$(shell which gocyclo)
HAVE_GOCOV:=$(shell which gocov)
GLIDE_VERSION='v0.12.2'
VERSION=$(patsubst "%",%,$(lastword $(shell grep 'const version' kamimai.go)))
COMMITISH=$(shell git rev-parse HEAD)
PROJECT_REPONAME=kamimai
PROJECT_USERNAME=eure
ARTIFACTS_DIR=artifacts

.PHONY: all ansible unit lint vet test

init: install-deps

build: install-deps

unit: lint vet cyclo build test
unit-report: lint vet cyclo build test-report

lint: golint
	@echo "go lint"
	@lint=`golint ./...`; \
	lint=`echo "$$lint" | grep -E -v -e ${LINTIGNOREDEPS}`; \
	echo "$$lint"; \
	if [ "$$lint" != "" ]; then exit 1; fi

vet:
	@echo "go vet"
	@go tool vet -all -structtags -shadow $(shell ls -d */ | grep -v "misc" | grep -v "vendor")

cyclo: gocyclo
	@echo "gocyclo -over 30"
	@gocyclo -over 30 ./core

test:
	@go test $(TARGET_ONLY_PKGS)

coverage: gocov
	@gocov test $(TARGET_ONLY_PKGS) | gocov report

test-report:
	@echo "Invoking test and coverage"
	@echo "" > coverage.txt; \
	for d in $(TARGET_ONLY_PKGS); do \
		go test -coverprofile=profile.out -covermode=atomic $$d || exit 1; \
		[ -f profile.out ] && cat profile.out >> coverage.txt && rm profile.out || true; \
	done

install-deps: glide
	@echo "Installing all dependencies"
	@PATH=$(INTERNAL_BIN):$(PATH) glide i

golint:
ifndef HAVE_GOLINT
	@echo "Installing linter"
	@go get -u github.com/golang/lint/golint
endif

gocyclo:
ifndef HAVE_GOCYCLO
	@echo "Installing gocyclo"
	@go get -u github.com/fzipp/gocyclo
endif

gocov:
ifndef HAVE_GOCOV
	@echo "Installing gocov"
	@go get -u github.com/axw/gocov/gocov
endif

glide:
ifndef HAVE_GLIDE
	@echo "Installing glide"
	@mkdir -p $(INTERNAL_BIN)
	@wget -q -O - https://github.com/Masterminds/glide/releases/download/$(GLIDE_VERSION)/glide-$(GLIDE_VERSION)-$(GOOS)-$(GOARCH).tar.gz | tar xvz
	@mv $(GOOS)-$(GOARCH)/glide $(INTERNAL_BIN)/glide
	@rm -rf $(GOOS)-$(GOARCH)
endif

ghr:
ifndef HAVE_GHR
	@echo "Installing ghr to upload binaries for release page"
	@go get -u github.com/tcnksm/ghr
endif

gox:
ifndef HAVE_GOX
	@echo "Installing gox to build binaries for Go cross compilation"
	@go get -u github.com/mitchellh/gox
endif

verify-github-token:
	@if [ -z "$$GITHUB_TOKEN" ]; then echo '$$GITHUB_TOKEN is required'; exit 1; fi

gox-build: gox
	@mkdir -p $(ARTIFACTS_DIR)/$(VERSION) && cd $(ARTIFACTS_DIR)/$(VERSION); \
		gox -ldflags="-s -w" github.com/eure/kamimai/cmd/kamimai

release: ghr verify-github-token gox-build
	@ghr -c $(COMMITISH) -u $(PROJECT_USERNAME) -r $(PROJECT_REPONAME) -t $$GITHUB_TOKEN \
		--replace $(VERSION) $(ARTIFACTS_DIR)/$(VERSION)
