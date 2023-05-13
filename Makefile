LDFLAGS = -s
LDFLAGS += -w
LDFLAGS += -X 'main.Version=$(shell git describe --tags --abbrev=0)'
LDFLAGS += -X 'main.Commit=$(shell git rev-list -1 HEAD)'
LDFLAGS += -extldflags "-static"'

clean:
	@if [ -e "./columbus-server" ];     then rm -rf "./columbus-server"     ; fi
	@if [ -e "./release/" ];      		then rm -rf "./release/"			; fi


build-prod: 
	go build -o columbus-server -tags netgo -ldflags="$(LDFLAGS)" .

build-dev:
	go build --race -o columbus-server .

build: build-prod

release: clean
	@mkdir release
	@go build -o release/columbus-server -tags netgo -ldflags="$(LDFLAGS)" .
	@cp server.conf release/
	@cp columbus-server.service release/
	@cd release/ && sha512sum * | gpg --local-user daniel@elmasy.com -o checksum.txt --clearsign
