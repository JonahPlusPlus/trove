BINARY_NAME=./out/trove.out

all: build run

build:
	go build -o ${BINARY_NAME} ./cmd/main.go

run:
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}
