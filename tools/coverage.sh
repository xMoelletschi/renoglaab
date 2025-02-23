#!/usr/bin/env bash

PKG_LIST=$(go list ./...)

mkdir -p "coverage"

go test -covermode=count -coverprofile "coverage/coverage.cov" ${PKG_LIST}

go tool cover -func="coverage/coverage.cov"

rm -rf "coverage"
