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
	Server.Router.StaticDirFS("assets", "资源文件", web.Assets, "assets")
	Server.Router.StaticFileFS("favicon.svg", "favicon", web.Favicon, "favicon.svg")
	// 首页 <index:path> 为路由参数，匹配任意路径，兼容单页应用
	Server.Router.StaticFileFS("<index:path>", "访问页面", web.IndexHtml, "index.html")
}
