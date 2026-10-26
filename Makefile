include Makefile.variables
include Makefile.precommit

DOCKER_REGISTRY ?= docker.io
IMAGE ?= bborbe/agent-code
ifeq ($(VERSION),)
	VERSION := $(shell git describe --tags `git rev-list --tags --max-count=1`)
endif

# Included AFTER VERSION is assigned: Makefile.docker's guard is evaluated at parse
# time, so including it above this block would check an empty VERSION, pass, and let
# the tag through to the shell lines that interpolate it.
include Makefile.docker

# Local dev only. Runtime config (kafka brokers, sentry DSN, stage vars) and the
# full deployment moved to the quant config repo — this is now a publish-only
# source repo (see CHANGELOG). Pass any flags/env locally as needed.
run:
	@go run -mod=mod main.go -v=2

deps:
	go install github.com/onsi/ginkgo/v2/ginkgo@v2.25.3
	sudo port install trivy

.PHONY: fix
fix:
	@for dir in $$(find `pwd` -type d -name vendor -prune -o -name go.mod -exec dirname "{}" \; | grep -v '^$$'); do \
		cd $${dir}; \
		echo "fix $${dir}"; \
		go get \
		github.com/bborbe/kv@latest \
		github.com/bborbe/memorykv@latest \
		github.com/bborbe/badgerkv@latest \
		github.com/bborbe/boltkv@latest \
		github.com/go-git/go-git/v5@latest \
		github.com/containerd/containerd@latest \
		golang.org/x/crypto@latest \
		golang.org/x/net@latest; \
	done
