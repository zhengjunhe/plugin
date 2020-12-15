#!/bin/sh
# proto生成命令，将pb.go文件生成到types/目录下, dplatformos_path支持引用dplatformos框架的proto文件
dplatformos_path=$(go list -f '{{.Dir}}' "github.com/33cn/dplatformos")
protoc --go_out=plugins=grpc:../types ./*.proto --proto_path=. --proto_path="${dplatformos_path}/types/proto/"
