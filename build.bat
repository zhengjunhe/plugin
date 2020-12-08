go env -w CGO_ENABLED=0
go build -o dplatform.exe
go build -o dplatform-cli.exe github.com/33cn/plugin/cli
