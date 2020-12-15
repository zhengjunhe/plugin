all:
	go get -v -u gitlab.33.cn/dplatformos/dplatformos
	go build -i -o tool github.com/33cn/dplatformos/cmd/tools
	./tool import --path "." --packname "github.com/33cn/plugin" --conf "" --out "plugin.toml"
