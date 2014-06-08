export GOPATH=$(shell pwd)

test:
	go build scheme
	go test -v eval
clean:
	rm -rf ./scheme
