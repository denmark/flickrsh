BINARY  := flickrsh
BUILD_DIR := bin

.PHONY: all build install run test vet fmt clean deps update-deps tidy

all: build

build:
	go build -o $(BUILD_DIR)/$(BINARY) .

install:
	go install .

run: build
	./$(BUILD_DIR)/$(BINARY)

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

clean:
	rm -rf $(BUILD_DIR)

deps:
	go mod download

tidy:
	go mod tidy

update-deps:
	go get -u ./...
	go mod tidy
