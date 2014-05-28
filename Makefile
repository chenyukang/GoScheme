export GOPATH=$(shell pwd)

test:
	go test -v eval
build:
	go build GoScheme

clean:
	rm -rf ./GoScheme
