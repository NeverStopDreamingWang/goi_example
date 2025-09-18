package mongodb

import (
	"github.com/NeverStopDreamingWang/goi"
	"goi_example/backend/database"
)

func init() {
	// 子路由
	mongodbRouter := database.DatabaseRouter.Include("mongodb/", "父路由")
	{
		mongodbRouter.Path("", "获取列表/创建", goi.ViewSet{GET: listView, POST: createView})
		// object_id 转换器将字符串转换为 primitive.ObjectID 类型
		mongodbRouter.Path("<object_id:pk>", "详情/更新/删除", goi.ViewSet{GET: retrieveView, PUT: updateView, DELETE: deleteView})
	}
}
