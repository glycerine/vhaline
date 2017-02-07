.phony: build test

nothing:

versionit:
	/bin/echo "package vhaline" > gitcommit.go
	/bin/echo "func init() { LAST_GIT_COMMIT_HASH = \"$(shell git rev-parse HEAD)\"; NEAREST_GIT_TAG= \"$(shell git describe --abbrev=0 --tags)\"; GIT_BRANCH=\"$(shell git rev-parse --abbrev-ref  HEAD)\"; GO_VERSION=\"$(shell go version)\";}" >> gitcommit.go

test:
	go test -v
