package goi_example

import (
	"goi_example/server/web"

	"github.com/NeverStopDreamingWang/goi"
)

var ApiRouter *goi.MetaRouter

func init() {
	Server.Router.StaticDir("static/", "静态目录", "static")
	ApiRouter = Server.Router.Include("api/", "API")
}

func InitIndexPage() {
	// 前端页面
	Server.Router.StaticFileFS("favicon.svg", "访问页面", web.Favicon, "favicon.svg")
	Server.Router.StaticDirFS("assets", "资源文件", web.Assets)
	Server.Router.Path("<index:path>", "访问页面", goi.ViewSet{GET: indexView})
}
