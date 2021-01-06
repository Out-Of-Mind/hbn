.PHONY: build
build:
	go build -v ./src/hbn.go

.DEFAULT_GOAL := build
