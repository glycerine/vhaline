.phony: build clean

build:
	cd ../../vhaline && make versionit
	CGGO_ENABLED=0 go build && cp -p ./vhaline $(GOPATH)/bin

clean:
	rm ./vhaline
