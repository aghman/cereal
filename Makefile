GO=/usr/local/go/bin/go

.DEFAULT_GOAL=build

.PHONY: setup build
setup:
	go mod vendor
	go mod tidy
	go mod vendor


.PHONY: build build_arm build_arm64 build_darwin_arm64

build: build_arm build_arm64 build_darwin_arm64

build_arm:
	GOOS=linux GOARCH=arm go build -o build/linux/arm/cereal main.go 
build_arm64:
	GOOS=linux GOARCH=arm64 go build -o build/linux/arm64/cereal main.go
build_darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -o build/darwin/arm64/cereal main.go