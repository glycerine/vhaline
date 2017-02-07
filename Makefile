.phony: build test

build:
	cd cmd/vhaline && make

test:
	cd vhaline && go test -v
