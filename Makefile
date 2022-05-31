# Project Setup
PROJECT_NAME := controller-util
PROJECT_REPO := github.com/presslabs/$(PROJECT_NAME)

PLATFORMS = linux_amd64 darwin_amd64

GO_SUBDIRS := pkg

GO111MODULE=on

GO_STATIC_PACKAGES = $(GO_PROJECT)/cmd/wp-operator
GO_LDFLAGS += -X $(PROJECT_REPO)/pkg/version.buildDate=$(BUILD_DATE) \
	       -X $(PROJECT_REPO)/pkg/version.gitVersion=$(VERSION) \
	       -X $(PROJECT_REPO)/pkg/version.gitCommit=$(GIT_COMMIT) \
	       -X $(PROJECT_REPO)/pkg/version.gitTreeState=$(GIT_TREE_STATE)

include build/makelib/common.mk
include build/makelib/golang.mk
include build/makelib/kubebuilder.mk
