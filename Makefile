# include the common make file
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
endif

# Copy githook scripts when execute makefile
COPY_GITHOOK:=$(shell cp -f githooks/* .git/hooks/)

# ==============================================================================
# Includes

include scripts/make-rules/tools.mk

# ==============================================================================

## tools: install dependent tools.
.PHONY: tools
tools:
	@$(MAKE) tools.install

.PHONY: lint
## lint: Check syntax and styling of go sources.
lint: tools.verify.golangci-lint
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
