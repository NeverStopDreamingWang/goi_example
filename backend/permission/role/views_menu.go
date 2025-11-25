package role

import (
	"net/http"

	"github.com/NeverStopDreamingWang/goi"
	"github.com/NeverStopDreamingWang/goi/db"
	"github.com/NeverStopDreamingWang/goi/db/sqlite3"
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
	var params roleMenuListValidParams
	var queryParams goi.Params
	var validationErr goi.ValidationError

	queryParams = request.QueryParams() // Query 传参
	validationErr = queryParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	var menu_list []*menuListModel
	sqlite3DB.SetModel(MenuModel{})
	err := sqlite3DB.Select(&menu_list)
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	for _, menu := range menu_list {
		sqlite3DB.SetModel(RoleMenuModel{})
		count, err := sqlite3DB.Where("role_id = ? and menu_id = ?", params.RoleId, menu.Id).Count()
		if err != nil {
			return goi.Data{
				Code:    http.StatusInternalServerError,
				Message: "查询数据库错误",
			}
		}
		menu.Checked = count != 0
	}
	// 获取树形结构
	menu_list = get_children_menu(menu_list)
	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: menu_list,
	}
}
