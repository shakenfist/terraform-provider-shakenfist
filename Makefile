# ShakenFist Terraform Provider makefile
#
# Suggestions taken from:
# https://github.com/terraform-providers/terraform-provider-aws/tree/master/aws
#
# For linter installation (in CI) see:
# 	https://golangci-lint.run/usage/install/#ci-installation
#


GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=provider
TEST?=./$(PKG_NAME)/...
TEST_COUNT?=1

default: build


build: fmtcheck
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

lint:
	@golangci-lint run ./$(PKG_NAME)/...

install-tools:
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint


.PHONY: build lint sweep test testacc fmt fmtcheck lint install-tools

