package user

import (
	"database/sql"
	"errors"
	"net/http"

	"goi_example/backend/permission/role"

	"github.com/NeverStopDreamingWang/goi"
	"github.com/NeverStopDreamingWang/goi/db"
	"github.com/NeverStopDreamingWang/goi/db/sqlite3"
)

// 参数验证
type listValidParams struct {
	Page         int64           `name:"page" type:"int" required:"true"`
	PageSize     int64           `name:"page_size" type:"int" required:"true"`
	Username     *string         `name:"username" type:"string"`
	Status       *UserStatusType `name:"status" type:"int"`
	RoleId       *int64          `name:"role_id" type:"int"`
	DepartmentId *int64          `name:"department_id" type:"int"`
}

type userList struct {
	Id            *int64          `json:"id"`
	Username      *string         `json:"username"`
	Email         *string         `json:"email"`
	Status        *UserStatusType `json:"status"`
	RoleId        *int64          `json:"role_id"`
	LastLoginTime *string         `json:"last_login_time"`
	CreateTime    *string         `json:"create_time"`
	UpdateTime    *string         `json:"update_time"`
	Role          *role.RoleModel `json:"role"`
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

	sqlite3DB.SetModel(UserModel{}) // 设置操作表

	if params.Username != nil {
		sqlite3DB = sqlite3DB.Where("`username` like ?", "%"+*params.Username+"%")
	}
	if params.Status != nil {
		sqlite3DB = sqlite3DB.Where("`status` = ?", params.Status)
	}
	if params.RoleId != nil {
		sqlite3DB = sqlite3DB.Where("`role_id` = ?", params.RoleId)
	}
	if params.DepartmentId != nil {
		sqlite3DB = sqlite3DB.Where("`department_id` = ?", params.DepartmentId)
	}
	total, total_page, err := sqlite3DB.Page(params.Page, params.PageSize)
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询用户失败",
			Results: err.Error(),
		}
	}

	var user_list []*userList
	err = sqlite3DB.Select(&user_list)
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询用户失败",
			Results: err.Error(),
		}
	}

	for _, user := range user_list {
		sqlite3DB.SetModel(role.RoleModel{})
		user.Role = &role.RoleModel{}
		err = sqlite3DB.Where("`id` = ?", user.RoleId).First(user.Role)
		if err != nil {
			continue
		}
	}

	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: map[string]any{
			"total":      total,
			"page":       params.Page,
			"total_page": total_page,
			"list":       user_list,
		},
	}
}

// 参数验证
type createValidParams struct {
	Username string         `name:"username" type:"string" required:"true"`
	Password string         `name:"password" type:"string" required:"true"`
	Email    string         `name:"email" type:"string" required:"true"`
	Status   UserStatusType `name:"status" type:"int" required:"true"`
	RoleId   int64          `name:"role_id" type:"int" required:"true"`
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

	user := &UserModel{
		Username: &params.Username,
		Password: &params.Password,
		Email:    &params.Email,
		Status:   &params.Status,
		RoleId:   &params.RoleId,
	}

	// 参数验证
	err := user.Validate()
	if err != nil {
		return goi.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Results: nil,
		}
	}
	// 添加
	err = user.Create()
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Results: nil,
		}
	}

	return goi.Data{
		Code:    http.StatusOK,
		Message: "添加用户",
		Results: user,
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

	user := &UserModel{}
	sqlite3DB.SetModel(UserModel{})
	err := sqlite3DB.Where("`id` = ?", pk).First(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return goi.Data{
				Code:    http.StatusBadRequest,
				Message: "用户不存在",
				Results: nil,
			}
		}
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询数据库错误",
			Results: err.Error(),
		}
	}

	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: user,
	}
}

// 参数验证
type updateValidParams struct {
	Username *string         `name:"username" type:"string"`
	Password *string         `name:"password" type:"string"`
	Email    *string         `name:"email" type:"string"`
	Status   *UserStatusType `name:"status" type:"int"`
	RoleId   *int64          `name:"role_id" type:"int"`
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

	instance := &UserModel{}
	sqlite3DB.SetModel(UserModel{})
	err := sqlite3DB.Where("`id` = ?", pk).First(instance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return goi.Data{
				Code:    http.StatusBadRequest,
				Message: "用户不存在",
				Results: nil,
			}
		}
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询数据库错误",
			Results: err.Error(),
		}
	}

	user := &UserModel{
		Id:       instance.Id,
		Username: params.Username,
		Password: params.Password,
		Email:    params.Email,
		Status:   params.Status,
		RoleId:   params.RoleId,
	}

	// 参数验证
	err = user.Validate()
	if err != nil {
		return goi.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Results: nil,
		}
	}
	// 更新
	err = instance.Update(user)
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Results: nil,
		}
	}
	return goi.Data{
		Code:    http.StatusOK,
		Message: "修改成功",
		Results: instance,
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

	instance := &UserModel{}
	sqlite3DB.SetModel(UserModel{})
	err := sqlite3DB.Where("`id` = ?", pk).First(instance)
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "用户不存在",
		}
	}
	err = instance.Delete()
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "删除失败",
		}
	}

	return goi.Data{
		Code:    http.StatusOK,
		Message: "删除成功",
		Results: nil,
	}
}
