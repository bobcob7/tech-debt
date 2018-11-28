GOOS   := linux
GOARCH := amd64

GOPATH := $(shell go env GOPATH)
PATH   := $(PATH):$(GOPATH)/bin
SHELL  := env PATH=$(PATH) /bin/bash

all: tech-debt

.PHONY: clean prerequisites

rice-box.go: client/dist/index.html
	rice embed-go

client/dist/index.html: client/src/main.js client/src/messages/messages_pb.js
	cd client && npm run build

tech-debt: server.go messages/messages_pb.go rice-box.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build

client/src/messages/messages_pb.js: prerequisites messages.proto
	mkdir -p client/src/messages
	protoc --js_out=import_style=commonjs,binary:client/src/messages/ messages.proto

messages/messages_pb.go: messages.proto
	mkdir -p messages
	go get -u github.com/golang/protobuf/protoc-gen-go
	protoc --go_out=messages/ messages.proto

clean:
	rm -rf tech-debt* client/dist
	rm -rf client/src/messages messages