all:
	go get -v -u gitlab.33.cn/dplatform/dplatform
	go build -i -o tool github.com/33cn/dplatform/cmd/tools
	./tool import --path "." --packname "github.com/33cn/plugin" --conf "" --out "plugin.toml"
