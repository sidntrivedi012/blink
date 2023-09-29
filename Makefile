.PHONY: build

BIN := blink.bin

build:
	go mod download
	go build -o ${BIN} ./cmd/blink/.

run:
	./${BIN}