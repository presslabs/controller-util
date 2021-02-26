KUBEBUILDER_VERSION ?= 2.3.2
GOLANGCI_LINT_VERSION := v1.37.1
BINDIR ?= $(PWD)/bin

GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH ?= amd64

PATH := $(BINDIR):$(PATH)
SHELL := env PATH=$(PATH) /bin/sh

all: lint

lint:
	GO111MODULE=on $(BINDIR)/golangci-lint run ./...

test:
	GO111MODULE=on KUBEBUILDER_ASSETS=$(BINDIR) ginkgo \
		--randomizeAllSpecs --randomizeSuites --failOnPending \
		--cover --coverprofile cover.out --trace --race \
		./...

dependencies:
	test -d $(BINDIR) || mkdir $(BINDIR)
	GOBIN=$(BINDIR) go get -u github.com/onsi/ginkgo/ginkgo
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- -b $(BINDIR) $(GOLANGCI_LINT_VERSION)
	curl -sL https://github.com/kubernetes-sigs/kubebuilder/releases/download/v$(KUBEBUILDER_VERSION)/kubebuilder_$(KUBEBUILDER_VERSION)_$(GOOS)_$(GOARCH).tar.gz | tar -zx -C $(BINDIR) --strip-components=2
