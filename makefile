build:
	cp openapi.yaml server/
	go build -o columbus-server -ldflags="-s -w" .
	rm server/openapi.yaml

build-dev:
	cp openapi.yaml server/
	go build --race -o columbus-server-dev .
	rm server/openapi.yaml	
