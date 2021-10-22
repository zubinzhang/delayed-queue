# ==============================================================================
# Makefile helper functions for tools
#

TOOLS ?= pre-commit golangci-lint

.PHONY: tools.install
tools.install: $(addprefix tools.verify., $(TOOLS))

.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: install.golangci-lint
install.golangci-lint:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.42.1

.PHONY: install.pre-commit
install.pre-commit:
	@curl https://pre-commit.com/install-local.py | python -
