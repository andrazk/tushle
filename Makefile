PROJECT=tushle

BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD | sed -E 's/[^a-zA-Z0-9]+/-/g')

LDFLAGS    := -ldflags " \
	-w -s \
	-X $(PROJECT)/cmd.revision=$(GIT_COMMIT) \
	-X $(PROJECT)/cmd.buildTime=$(BUILD_TIME) \
"

BUILDER_PATH  := resources/docker
BUILDER_IMAGE := $(PROJECT)/builder:0.1

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show help
	@echo "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:"
	@grep -E '^[a-zA-Z_/%\-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

deploy: build
	scp "${CURDIR}/dist/tushle" "tushle-prod:/root/tushle"

.PHONY: binary-osx
binary-osx: ## build executable for macOS
	./scripts/build/osx

docker/build/builder: ## Build builder image
	docker build -t $(BUILDER_IMAGE) -f $(BUILDER_PATH)/Dockerfile $(BUILDER_PATH)

builder/%:: ## Runs make target in builder container
	docker run -it \
		-v "$(shell pwd)":/go/src/$(PROJECT) \
		-w /go/src/$(PROJECT) \
		$(BUILDER_IMAGE) make $*;

console: ## Runs bash shell, i.e. builder/console
	@bash

test/static: ## Test for unused code
	staticcheck tushle/...
