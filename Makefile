PKG_LIST := $(shell go list ./...)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)

.PHONY: build
build:
	go build -race -o bin/renoglaab cmd/renoglaab/main.go

.PHONY: run
run: build
	./bin/renoglaab

.PHONY: lint-test
lint-test:
	golangci-lint run -v

.PHONY: staticcheck-test
staticcheck-test:
	staticcheck ${PKG_LIST}

.PHONY: fmt-test
fmt-test:
	go fmt ${PKG_LIST}

.PHONY: vet-test
vet-test:
	go vet ${PKG_LIST}

.PHONY: unit-test
unit-test:
	go test -v ${PKG_LIST} -count=1 -timeout=10s

.PHONY: race-test
race-test:
	go test -race -short ${PKG_LIST}  -count=1 -timeout=10s

.PHONY: gosec-test
gosec-test:
	gosec ${PKG_LIST}

.PHONY: coverage-test
coverage-test:
	./tools/coverage.sh

.PHONY: benchmark-test
benchmark-test:
	go test -bench=. -benchmem ${PKG_LIST}

# Just for using locally to see the coverage report via html format.
.PHONY: coverage-html
coverage-html:
	mkdir -p "coverage"
	go test -covermode=count -coverprofile "coverage/coverage.cov" ${PKG_LIST}
	go tool cover -func="coverage/coverage.cov"
	go tool cover -html="coverage/coverage.cov" -o "bin/coverage.html"
	open "bin/coverage.html"
	rm -rf "coverage"

.PHONY: gotestsum
gotestsum:
	gotestsum --watch -- --count=1 --timeout=5s

.PHONY: local-test-moelletschi
local-test-moelletschi:
	LOG_LEVEL=debug CONFIG_PATH=bin/config.js AUTHOR_USERNAME=xMoelletschi go run cmd/renoglaab/main.go

.PHONY: local-test-moelletschi-branch
local-test-moelletschi-branch:
	LOG_LEVEL=debug CONFIG_PATH=bin/config.js AUTHOR_USERNAME=xMoelletschi ALLOWED_BRANCH_REGEX="renovate/(master|main)-automerge" go run cmd/renoglaab/main.go
