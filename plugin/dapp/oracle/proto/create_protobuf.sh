#!/bin/sh

dplatform_path=$(go list -f '{{.Dir}}' "github.com/33cn/dplatform")
protoc --go_out=plugins=grpc:../types ./*.proto --proto_path=. --proto_path="${dplatform_path}/types/proto/"
