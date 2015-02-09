export GOPATH := $(shell pwd)

all:
	make deps
	make build

deps:
	go get github.com/go-martini/martini
	go get github.com/vanng822/imgscale/imagick
	go get github.com/vanng822/imgscale/imgscale

build:
	go build -o bin/imgscale
	
install:
	go install

clean:
	rm -r pkg/