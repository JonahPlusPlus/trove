SERVER_OUT=./out/trove.out
WASM_OUT=./static/dashboard/index.wasm

all: templates scss wasm server run

server:
	go build -o ${SERVER_OUT} ./cmd/server/main.go

wasm:
	GOOS=js GOARCH=wasm go build -o ${WASM_OUT} ./cmd/dashboard/main.go

scss:
	sass static/scss:static/css --style=compressed

run:
	./${SERVER_OUT}

clean:
	go clean
	rm ${SERVER_OUT}

install: install_quicktemplate certificate

certificate:
	./generate_certificate.sh

install_quicktemplate:
	go install github.com/valyala/quicktemplate/qtc

templates:
	qtc -dir=templates

boot:
	docker compose up
