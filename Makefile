# ShakenFist Terraform Provider makefile
#
# Suggestions taken from:
# https://github.com/terraform-providers/terraform-provider-aws/tree/master/aws
#
# For linter installation (in CI) see:
# 	https://golangci-lint.run/usage/install/#ci-installation
#
# Goreleaser installation:
#	https://github.com/goreleaser/goreleaser/releases/latest
#

# Settings
NAME = shakenfist
PKG_NAME = provider
GPG_KEY = CEDCE5DB21914905D930A42CF31CB8C24C064BC3

# Setup
BINARY = terraform-provider-${NAME}
VERSION = $(shell git describe --abbrev=0 --tags)
VERSION_NO_V = $(VERSION:v%=%)

# Test and lint
GOFMT_FILES ?= $(shell find . -name '*.go' |grep -v vendor)
TEST ?= ./$(PKG_NAME)/...
TEST_COUNT ?= 1


default: build

build: fmtcheck $(GOFMT_FILES)
	@echo "==> Building..."
	go build .

# Acceptance tests
testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v -count $(TEST_COUNT) -parallel 20 $(TESTARGS) -timeout 120m

# Unit tests
test: fmtcheck
	go test $(TEST) $(TESTARGS) -v -timeout=120s -parallel=4

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/check_files_gofmt.sh'"

install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/shakenfist/${NAME}/${VERSION_NO_V}/linux_amd64
	cp ${BINARY} ~/.terraform.d/plugins/registry.terraform.io/shakenfist/${NAME}/${VERSION_NO_V}/linux_amd64

lint:
	@golangci-lint run ./$(PKG_NAME)/...

install-tools:
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint

release:
	GPG_FINGERPRINT=$(GPG_KEY) goreleaser --rm-dist --skip-publish

.PHONY: build lint sweep test testacc fmt fmtcheck lint install-tools release
