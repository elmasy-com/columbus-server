LDFLAGS = -s
LDFLAGS += -w
LDFLAGS += -X 'main.Version=$(shell git describe --tags --abbrev=0)'
LDFLAGS += -X 'main.Commit=$(shell git rev-list -1 HEAD)'
LDFLAGS += -extldflags "-static"'

clean:
	@if [ -e "./columbus-server" ];     then rm -rf "./columbus-server"     ; fi
	@if [ -e "./columbus-frontend/" ];  then rm -rf "./columbus-frontend"   ; fi
	@if [ -e "./server/static/" ];      then rm -rf "./server/static/"   	; fi
	@if [ -e "./release/" ];      		then rm -rf "./release/"			; fi


static: clean
	@git clone git@github.com:elmasy-com/columbus-frontend.git
	@cd columbus-frontend && npm install && npm run build
	@mv columbus-frontend/dist server/static
	@cp openapi.yaml server/static/
	@rm -rf columbus-frontend

build-prod: static
	go build -o columbus-server -tags netgo -ldflags="$(LDFLAGS)" .

build-dev: static
	go build --race -o columbus-server .

build: build-prod

release: static
	@mkdir release
	@go build -o release/columbus-server -tags netgo -ldflags="$(LDFLAGS)" .
	@cp server.conf release/
	@cp columbus-server.service release/
	@cd release/ && sha512sum * | gpg --local-user daniel@elmasy.com -o checksum.txt --clearsign
