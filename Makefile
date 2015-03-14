export GOPATH := $(shell pwd)

all:
	make deps
	make build

deps:
	go get -u github.com/go-martini/martini
	go get -u github.com/vanng822/imgscale/imgscale
	go get -u github.com/vanng822/imgscale/provider/http

build:
	go build -o bin/imgscale
	
install:
	go install

clean:
	rm -r pkg/