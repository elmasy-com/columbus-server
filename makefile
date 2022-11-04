LDFLAGS = -s
LDFLAGS += -w
LDFLAGS += -X 'main.Version=$(shell git tag)'
LDFLAGS += -X 'main.Commit=$(shell git rev-list -1 HEAD)'

build:
	cp openapi.yaml server/
	go build -o columbus-server -ldflags="$(LDFLAGS)" .
	rm server/openapi.yaml

build-dev:
	cp openapi.yaml server/
	go build --race -o columbus-server-dev .
	rm server/openapi.yaml	
