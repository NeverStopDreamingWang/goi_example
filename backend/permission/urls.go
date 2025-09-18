package permission

import (
	"goi_example/backend/goi_example"

	"github.com/NeverStopDreamingWang/goi"
)

var PermissionRouter *goi.MetaRouter

func init() {
	// 子路由
	PermissionRouter = goi_example.ApiRouter.Include("permission/", "权限")
}
