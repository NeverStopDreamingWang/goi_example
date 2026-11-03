package role

import (
	"net/http"

	"github.com/NeverStopDreamingWang/goi/v3"
	"github.com/NeverStopDreamingWang/goi/v3/db"
	"github.com/NeverStopDreamingWang/goi/v3/db/sqlite3"
	"github.com/NeverStopDreamingWang/goi/v3/response"
)

// 参数验证
type roleMenuListValidParams struct {
	RoleId int `name:"role_id" type:"int" required:"true"`
}

type menuListModel struct {
	Id       *int64           `json:"id"`
	ParentId *int64           `json:"parent_id"`
	Name     *string          `json:"name"`
	Icon     *string          `json:"icon"`
	Path     *string          `json:"path"`
	Checked  bool             `json:"checked"`
	Children []*menuListModel `json:"children"`
}

func roleMenuListView(request *goi.Request) any {
	ctx := request.Object.Context()
	var params roleMenuListValidParams
	var queryParams goi.Params
	var validationErr goi.ValidationError

	queryParams = request.QueryParams() // Query 传参
	validationErr = queryParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	var menuList []*menuListModel
	sqlite3DB.SetModel(MenuModel{})
	err := sqlite3DB.Select(ctx, &menuList)
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	for _, menu := range menuList {
		sqlite3DB.SetModel(RoleMenuModel{})
		count, err := sqlite3DB.Where("role_id = ? and menu_id = ?", params.RoleId, menu.Id).Count(ctx)
		if err != nil {
			return response.Data{
				Code:    http.StatusInternalServerError,
				Message: "查询数据库错误",
			}
		}
		menu.Checked = count != 0
	}
	// 获取树形结构
	menuList = get_children_menu(menuList)
	return response.Data{
		Code:    http.StatusOK,
		Message: "",
		Data:    menuList,
	}
}
