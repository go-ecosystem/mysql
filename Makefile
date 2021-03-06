SHELL := /bin/bash
BASEDIR = $(shell pwd)

all: fmt mod lint test
fmt:
	gofmt -w .
mod:
	go mod tidy
lint:
	golangci-lint run
.PHONY: test
test: mod
	sh scripts/test.sh
.PHONY: mysql
mysql:
	sh scripts/mysql.sh
help:
	@echo "fmt - format the source code"
	@echo "mod - go mod tidy"
	@echo "lint - run golangci-lint"
	@echo "test - unit test"
	@echo "mysql - launch a docker mysql"