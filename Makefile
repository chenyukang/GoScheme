export GOPATH=$(shell pwd)


all:
	go build scheme

test:
	go build scheme
	go test -v eval
clean:
	rm -rf ./scheme
