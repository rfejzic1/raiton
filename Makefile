BINARY=raiton
BUILD_DIR=./build
SRC=cmd/cli/main.go

all: build

.PHONY: build
build:
	@go build -o $(BUILD_DIR)/${BINARY} $(SRC)

.PHONY: run
run: build
	@${BUILD_DIR}/$(BINARY)

.PHONY: testall
testall:
	@go test ./...

.PHONY: clean
clean:
	@go clean
	@rm -rf $(BUILD_DIR)

