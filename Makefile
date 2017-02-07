.phony: build test

build:
	cd vhaline && make && cd ../cmd/vhaline && make

test:
	cd vhaline && go test -v
