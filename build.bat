go env -w CGO_ENABLED=0
go build -o dplatformos.exe
go build -o dplatformos-cli.exe github.com/33cn/plugin/cli
