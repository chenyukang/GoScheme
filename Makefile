export GOPATH=$(shell pwd)

test:
	go test -v eval
	go build scheme

clean:
	rm -rf ./scheme
