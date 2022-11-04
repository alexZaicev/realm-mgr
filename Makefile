GO   := go
BUF  := buf
WIRE := wire

DIR_CMD      := ./cmd
DIR_DIST     := ./dist
DIR_MOCKS    := ./mocks
DIR_INTERNAL := ./internal

MOCKERY      := mockery
MOCKERY_ARGS := --all --keeptree --dir $(DIR_INTERNAL)

VERSION = $(shell (git describe --long --tags --match 'v[0-9]*' || echo v0.0.0) | cut -c2-)
COMMIT  = $(shell git rev-parse --short HEAD)
LDFLAGS = -X main.Version=$(VERSION)

INTERNAL_NON_TEST_GO_FILES = $(shell find $(DIR_INTERNAL) -type f -name '*.go' -not -name '*_test.go')

# -----------------------------------------------------------------
# Database build targets
# -----------------------------------------------------------------

.PHONY: init_dev_db
init_dev_db:
	./build/tools/database/create_dev_db.sh $(VERSION)

# -----------------------------------------------------------------
# Service build targets
# -----------------------------------------------------------------

.PHONY: vendor
vendor:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: checks
checks: wire mocks fmt

.PHONY: dirty
dirty:
	if [ $$(git status --porcelain | wc -l) -ne "0" ]; then \
		echo "Missing / modified files:"; \
		git status --porcelain; \
		echo; \
		echo "Diff of changed files:"; \
		git diff; \
		exit 1; \
	fi

.PHONY: unit
unit: vendor
	go test $(DIR_INTERNAL)/... $(DIR_CMD)/... \
		-cover \
		-coverprofile=coverage.out \
		-count=1
	@cat coverage.out | \
		awk 'BEGIN {cov=0; stat=0;} $$3!="" { cov+=($$3==1?$$2:0); stat+=$$2; } \
    	END {printf("Total coverage: %.2f%% of statements\n", (cov/stat)*100);}'
	go tool cover -html=coverage.out -o coverage.html

.PHONY: build
build: vendor
	@mkdir -p $(DIR_DIST)
	$(GO) build -ldflags "$(LDFLAGS)" -o $(DIR_DIST)/ $(DIR_CMD)/...

.PHONY: lint
lint: golint #helmlint

.PHONY: golint
golint:
	golangci-lint run --concurrency=2 --timeout=30m --max-issues-per-linter 0 --max-same-issues 0

.PHONY: wire
wire:
	$(WIRE) $(DIR_CMD)/realm-mgr-grpc

.PHONY: mocks
mocks: $(INTERNAL_NON_TEST_GO_FILES)
	rm -rf $(DIR_MOCKS)_maketemp/
	@# Mockery returns error code 0 on these errors but produces incorrect output
	if $(MOCKERY) $(MOCKERY_ARGS) --output $(DIR_MOCKS)_maketemp 2>&1 | grep ERR; then \
		rm -rf $(DIR_MOCKS)_maketemp; \
		exit 1; \
	fi
	rm -rf $(DIR_MOCKS)/
	mv $(DIR_MOCKS)_maketemp $(DIR_MOCKS)

## fmt: format the code
.PHONY: fmt
fmt:
	gofmt -s -w -e $(DIR_CMD) $(DIR_INTERNAL)
	gci write \
		-s Standard \
		-s Default \
		-s 'Prefix(github.com)' \
		-s 'Prefix(github.com/alexZaicev/realm-mgr)' \
		$(DIR_CMD) $(DIR_INTERNAL)
	goimports -local github.hpe.com -w $(DIR_CMD) $(DIR_INTERNAL)

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
