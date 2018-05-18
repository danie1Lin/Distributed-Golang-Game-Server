#!/bin/bash
protoc -I. --csharp_out . --grpc_out . message.proto --plugin=protoc-gen-grpc=/usr/sbin/grpc_csharp_plugin
protoc --go_out=plugins=grpc:./ message.proto
zip -r message.zip Message* message.proto

