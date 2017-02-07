.phony: build test

build:
	cd vhaline && go install && cd ../cmd/vhaline && go install

test:
	cd vhaline && go test -v
