package database

import (
	"goi_example/backend/goi_example"

	"github.com/NeverStopDreamingWang/goi"
)

var DatabaseRouter *goi.MetaRouter

func init() {
	// 子路由
	DatabaseRouter = goi_example.ApiRouter.Include("database/", "数据库")
}
