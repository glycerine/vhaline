.phony: build test

build:
	cd cmd/vhaline && go install && go build

test:
	cd vhaline && go test -v
