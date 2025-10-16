package goi_example

import (
	"goi_example/server/web"

	"github.com/NeverStopDreamingWang/goi"
)

var ApiRouter *goi.MetaRouter

func init() {
	// 前端页面
	Server.Router.StaticFileFS("", "访问页面", web.IndexHtml, "index.html")
	Server.Router.StaticFileFS("favicon.svg", "favicon", web.Favicon, "favicon.svg")
	Server.Router.StaticDirFS("assets", "资源文件", web.Assets, "assets")

	// 未匹配路由，返回首页兼容单页应用
	Server.Router.NoRoute(goi.ViewSet{
		HEAD: goi.StaticFileFSView(web.IndexHtml, "index.html"),
		GET:  goi.StaticFileFSView(web.IndexHtml, "index.html"),
	})

	Server.Router.StaticDir("static/", "静态目录", "static")
	ApiRouter = Server.Router.Include("api/", "API")
}
