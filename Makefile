THIS_FILE := $(lastword $(MAKEFILE_LIST))
SHELL=/bin/bash

all:
	# api service actions
	gridder-wire

	gridder-build_linux_amd64
	gridder-build_darwin_arm64
	gridder-build

###### VARIABLES ######
CMD_BUILD_DEV=go build -race -o ./bin/$(1)_dev ./cmd/$(1)/
CMD_BUILD_PRD=GOOS=$(1) GOARCH=$(2) go build -o ./bin/$(3) ./cmd/$(3)/

####### TARGETS #######
gridder-wire:
	@wire gen $(APP_API)

gridder-build_linux_amd64:
	@$(call CMD_BUILD,linux,amd64,api)

gridder-build_darwin_arm64:
	@$(call CMD_BUILD,darwin,arm64,api)

gridder-build:
	@$(call CMD_BUILD_DEV,gridder)
