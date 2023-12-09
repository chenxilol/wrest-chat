package main

import (
	"embed"

	"github.com/opentdp/wechat-rest/args"
	"github.com/opentdp/wechat-rest/httpd"
)

//go:embed public
var efs embed.FS

func main() {

	args.Efs = &efs
	args.NewConfig().Init()

	httpd.Server()

}
