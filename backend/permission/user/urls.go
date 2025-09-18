package user

import (
	"goi_example/backend/permission"

	"github.com/NeverStopDreamingWang/goi"
)

func init() {
	// 子路由
	userRouter := permission.PermissionRouter.Include("user/", "用户管理")
	{
		userRouter.Path("", "用户查看创建", goi.ViewSet{GET: listView, POST: createView})
		userRouter.Path("<int:pk>", "用户详情修改删除", goi.ViewSet{GET: retrieveView, PUT: updateView, DELETE: deleteView})

		userRouter.Path("profile", "用户详情修改删除", goi.ViewSet{GET: profileRetrieveView, PUT: profileUpdateView})
	}
}
