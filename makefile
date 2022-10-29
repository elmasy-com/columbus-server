build:
	cp openapi.yaml server/
	go build .
	rm server/openapi.yaml