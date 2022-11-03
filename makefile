build:
	cp openapi.yaml server/
	go build -o columbus-server .
	rm server/openapi.yaml

devbuild:
	cp openapi.yaml server/
	go build --race -o columbus-server-dev .
	rm server/openapi.yaml	
