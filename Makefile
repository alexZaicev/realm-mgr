GO  := go
BUF := buf

DIR_CMD  := ./cmd
DIR_DIST := ./dist

VERSION = $(shell (git describe --long --tags --match 'v[0-9]*' || echo v0.0.0) | cut -c2-)
COMMIT  = $(shell git rev-parse --short HEAD)
LDFLAGS = -X main.Version=$(VERSION)

# -----------------------------------------------------------------
# Service build targets
# -----------------------------------------------------------------

.PHONY: vendor
vendor:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: build
build: vendor
	@mkdir -p $(DIST_DIR)
	$(GO) build -ldflags "$(LDFLAGS)" -o $(DIR_DIST) $(DIR_CMD)/...

# -----------------------------------------------------------------
# Proto build targets
# -----------------------------------------------------------------

.PHONY: buf-checks
buf-checks: buf-lint buf-break buf-gen

.PHONY: buf-gen
buf-gen: buf-gen-code proto-mocks

.PHONY: buf-gen-code
buf-gen-code:
	@echo "Generating go code from proto files"
	$(BUF) generate $(shell find -path "./proto/proto/*" -prune -printf '--path "%p" '| sort) -o ./proto

.PHONY: buf-lint
buf-lint:
	@echo "Running lint check"
	$(BUF) lint

.PHONY: buf-break
buf-break:
	@echo "Running breaking changes check"
# Set the input to check against conditionally to get it working in CI
ifeq ("$(GITHUB_TOKEN)","")
	$(BUF) breaking --against .git\#branch=master
else
	$(BUF) breaking --against "https://x-access-token@${GITHUB_TOKEN}@github.com/alexZaicev/realm-mgr.git"
endif

.PHONY: proto-mocks
proto-mocks:
	for module in $(shell find proto/go -type d -regex ".*/v[0-9]*" | sort);\
	do\
		rm -rf $$module/mocks;\
		mkdir $$module/mocks;\
		GOFLAGS="" mockery --all --recursive --dir $$module --output "$$module/mocks";\
	done
