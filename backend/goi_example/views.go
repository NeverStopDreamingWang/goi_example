package goi_example

import (
	"github.com/NeverStopDreamingWang/goi"
	"goi_example/server/web"
)

func indexView(request *goi.Request) interface{} {
	request.Object.URL.Path = "/"
	return web.IndexHtml
}
