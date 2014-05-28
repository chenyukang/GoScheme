export GOPATH=$(shell pwd)

test:
	go test -v ./eval
build:
	go build scheme.go

clean:
	rm -rf ./scheme
