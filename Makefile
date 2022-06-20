SERVER_OUT=./out/trove.out
WASM_OUT=./static/dashboard/index.wasm

all: generate build_wasm build_server run

build_server:
	go build -o ${SERVER_OUT} ./cmd/server/main.go

build_wasm:
	GOOS=js GOARCH=wasm go build -o ${WASM_OUT} ./cmd/dashboard/main.go

run:
	./${SERVER_OUT}

clean:
	go clean
	rm ${SERVER_OUT}

certificate:
	./generate_certificate.sh

generate:
	go generate

boot:
	docker compose up