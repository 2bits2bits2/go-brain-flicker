BINARY_NAME=myapp
VERSION=1.0.0
BUILD_DIR=build

# Ensure build directory exists
$(shell mkdir -p ${BUILD_DIR})

.PHONY: clean build build-optim

clean:
	rm -rf ${BUILD_DIR}/*

# Regular builds
build: clean
	# Linux build
	GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 ./...
	# Windows build
	GOOS=windows GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe ./...

# Optimized (smaller) builds
build-optim: clean
	# Linux optimized build
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64-min ./...
	# Windows optimized build
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64-min.exe ./...

# Build all versions
all: build build-optim