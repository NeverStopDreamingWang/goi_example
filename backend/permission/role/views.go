package role

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/NeverStopDreamingWang/goi/v2"
	"github.com/NeverStopDreamingWang/goi/v2/db"
	"github.com/NeverStopDreamingWang/goi/v2/db/sqlite3"
	"github.com/NeverStopDreamingWang/goi/v2/response"
)

// 参数验证
type listValidParams struct {
	Page     int64  `name:"page" type:"int" required:"true"`
	PageSize int64  `name:"page_size" type:"int" required:"true"`
	Name     string `name:"name" type:"string"`
}

func listView(request *goi.Request) any {
	var params listValidParams
	var queryParams goi.Params
	var validationErr goi.ValidationError

	queryParams = request.QueryParams()
	validationErr = queryParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	sqlite3DB.SetModel(RoleModel{}) // 设置操作表

	if params.Name != "" {
		sqlite3DB = sqlite3DB.Where("`name` like ?", "%"+params.Name+"%")
	}
	total, total_page, err := sqlite3DB.Page(params.Page, params.PageSize)
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询角色失败",
			Data:    err.Error(),
		}
	}

	var roleList []RoleModel
	err = sqlite3DB.Select(&roleList)
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询角色失败",
			Data:    err.Error(),
		}
	}
	return response.Data{
		Code:    http.StatusOK,
		Message: "ok",
		Data: map[string]any{
			"total":      total,
			"page":       params.Page,
			"total_page": total_page,
			"list":       roleList,
		},
	}
}

// 参数验证
type createValidParams struct {
	Name     string   `name:"name" type:"string" required:"true"`
	MenuList []*int64 `name:"menu_list" type:"slice" required:"true"`
}

func createView(request *goi.Request) any {
	var params createValidParams
	var bodyParams goi.Params
	var validationErr goi.ValidationError

	bodyParams = request.BodyParams() // Body 传参
	validationErr = bodyParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	role := &RoleModel{
		Name:     &params.Name,
		MenuList: params.MenuList,
	}
	// 参数验证
	err := role.Validate()
	if err != nil {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}
	// 添加
	err = role.Create()
	if err != nil {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return response.Data{
		Code:    http.StatusOK,
		Message: "",
		Data:    role,
	}
}

func retrieveView(request *goi.Request) any {
	var pk int64
	var validationErr goi.ValidationError
	validationErr = request.PathParams.Get("pk", &pk) // 路由传参
	if validationErr != nil {
		return validationErr.Response()
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	role := &RoleModel{}
	sqlite3DB.SetModel(RoleModel{})
	err := sqlite3DB.Where("`id` = ?", pk).First(role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.Data{
				Code:    http.StatusBadRequest,
				Message: "角色不存在",
				Data:    nil,
			}
		}
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询角色失败",
			Data:    err.Error(),
		}
	}

	return response.Data{
		Code:    http.StatusOK,
		Message: "",
		Data:    role,
	}
}

// 参数验证
type updateValidParams struct {
	Name     string   `name:"name" type:"string" required:"true"`
	MenuList []*int64 `name:"menu_list" type:"slice"`
}

func updateView(request *goi.Request) any {
	var pk int64
	var params updateValidParams
	var bodyParams goi.Params
	var validationErr goi.ValidationError

	validationErr = request.PathParams.Get("pk", &pk) // 路由传参
	if validationErr != nil {
		return validationErr.Response()
	}

	bodyParams = request.BodyParams() // Body 传参
	validationErr = bodyParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	instance := &RoleModel{}
	sqlite3DB.SetModel(RoleModel{})
	err := sqlite3DB.Where("`id` = ?", pk).First(instance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.Data{
				Code:    http.StatusBadRequest,
				Message: "角色不存在",
				Data:    nil,
			}
		}
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询数据库错误",
			Data:    err.Error(),
		}
	}

	role := &RoleModel{
		Id:       instance.Id,
		Name:     &params.Name,
		MenuList: params.MenuList,
	}
	// 参数验证
	err = role.Validate()
	if err != nil {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}
	// 修改
	err = instance.Update(role)
	if err != nil {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return response.Data{
		Code:    http.StatusOK,
		Message: "修改成功",
		Data:    instance,
	}
}

func deleteView(request *goi.Request) any {
	var pk int64
	var validationErr goi.ValidationError

	validationErr = request.PathParams.Get("pk", &pk) // 路由传参
	if validationErr != nil {
		return validationErr.Response()
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	instance := &RoleModel{}
	sqlite3DB.SetModel(RoleModel{})
	err := sqlite3DB.Where("`id` = ?", pk).First(instance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.Data{
				Code:    http.StatusBadRequest,
				Message: "角色不存在",
				Data:    nil,
			}
		}
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询数据库错误",
			Data:    err.Error(),
		}
	}
	err = instance.Delete()
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "删除失败",
			Data:    err.Error(),
		}
	}

	return response.Data{
		Code:    http.StatusOK,
		Message: "删除成功",
		Data:    nil,
	}
}

// 参数验证
type selectValidParams struct {
	Name *string `name:"name" type:"string"`
}

type RoleSelect struct {
	Id   *int64  `json:"value"`
	Name *string `json:"label"`
}

func allView(request *goi.Request) any {
	var params selectValidParams
	var queryParams goi.Params
	var validationErr goi.ValidationError

	queryParams = request.QueryParams()
	validationErr = queryParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	sqlite3DB.SetModel(RoleModel{}) // 设置操作表

	if params.Name != nil {
		sqlite3DB.Where("`name` like ?", "%"+*params.Name+"%")
	}

	var roleList []RoleSelect
	err := sqlite3DB.Select(&roleList)
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询失败",
			Data:    nil,
		}
	}
	return response.Data{
		Code:    http.StatusOK,
		Message: "",
		Data:    roleList,
	}
}
