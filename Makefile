.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%H%M%S')

APP = mag
SERVER_BIN = ./bin/${APP}
RELEASE_ROOT = build
RELEASE_SERVER = build/${APP}

all: start

build:
	@go build -ldflags "-w -s" -o $(SERVER_BIN) .

start:
	go run main.go web -c ./conf/config.toml -r ./conf/casbin.conf -m ./conf/menu.yaml -w ./web/dist

swagger:
	swag init --generalInfo ./server/swagger --output ./server/swagger

wire:
	wire gen ./server/provider

test:
	@go test -v $(shell go list ./...)

clean:
	rm -fr build ${SERVER_BIN}

pack: build
	rm -fr $${RELEASE_ROOT} && mkdir -p ${RELEASE_SERVER}
	cp -r ${SERVER_BIN} conf ${RELEASE_SERVER}
	cd ${RELEASE_ROOT} && tar -cvf ${APP}.tar ${APP} && rm -fr ${APP}