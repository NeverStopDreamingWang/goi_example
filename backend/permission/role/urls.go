package role

import (
	"goi_example/backend/permission"

	"github.com/NeverStopDreamingWang/goi"
)

func init() {
	// 子路由
	roleRouter := permission.PermissionRouter.Include("role/", "角色管理")
	{
		// 角色
		roleRouter.Path("", "角色查看创建", goi.ViewSet{GET: listView, POST: createView})
		roleRouter.Path("<int:pk>", "角色详情修改删除", goi.ViewSet{GET: retrieveView, PUT: updateView, DELETE: deleteView})

		// 菜单
		roleRouter.Path("menu", "角色菜单列表", goi.ViewSet{GET: roleMenuListView})

		roleRouter.Path("select", "角色下拉菜单", goi.ViewSet{GET: allView})
	}
}
