package mysql

import (
	"goi_example/backend/database"

	"github.com/NeverStopDreamingWang/goi"
)

func init() {
	// 子路由
	mysqlRouter := database.DatabaseRouter.Include("mysql/", "父路由")
	{
		// 注册一个路径
		mysqlRouter.Path("", "获取列表/创建", goi.ViewSet{GET: listView, POST: createView})
		mysqlRouter.Path("<int:pk>", "详情/更新/删除", goi.ViewSet{GET: retrieveView, PUT: updateView, DELETE: deleteView})

	}
}
