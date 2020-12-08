#!/bin/sh
# proto生成命令，将pb.go文件生成到types/目录下, dplatform_path支持引用dplatform框架的proto文件
dplatform_path=$(go list -f '{{.Dir}}' "github.com/33cn/dplatform")
protoc --go_out=plugins=grpc:../types ./*.proto --proto_path=. --proto_path="${dplatform_path}/types/proto/"
