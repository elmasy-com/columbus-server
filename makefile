LDFLAGS = -s
LDFLAGS += -w
LDFLAGS += -X 'main.Version=$(shell git describe --tags --abbrev=0)'
LDFLAGS += -X 'main.Commit=$(shell git rev-list -1 HEAD)'

build:
	cp openapi.yaml server/
	go build -o columbus-server -ldflags="$(LDFLAGS)" .
	rm server/openapi.yaml

build-dev:
	cp openapi.yaml server/
	go build --race -o columbus-server-dev .
	rm server/openapi.yaml	
