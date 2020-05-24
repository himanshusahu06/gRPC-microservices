#!/bin/sh
#go get github.com/golang/protobuf/protoc-gen-go
protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.