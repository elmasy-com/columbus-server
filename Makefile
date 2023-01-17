LDFLAGS = -s
LDFLAGS += -w
LDFLAGS += -X 'main.Version=$(shell git describe --tags --abbrev=0)'
LDFLAGS += -X 'main.Commit=$(shell git rev-list -1 HEAD)'

clean:
	@if [ -e "./columbus-server" ];     then rm -rf columbus-server     ; fi
	@if [ -e "./columbus-server-dev" ]; then rm -rf columbus-server-dev ; fi
	@if [ -e "./columbus-frontend/" ];  then rm -rf columbus-frontend   ; fi
	@if [ -e "./server/static/" ];      then rm -rf "server/static"     ; fi

static: clean
	@git clone git@github.com:elmasy-com/columbus-frontend.git
	@cd columbus-frontend && npm install && npm run build
	@mv columbus-frontend/dist server/static
	@cp openapi.yaml server/static/
	@rm -rf columbus-frontend


build-prod: static
	go build -o columbus-server -ldflags="$(LDFLAGS)" .

build-dev: static
	go build --race -o columbus-server .

build: static build-prod clean