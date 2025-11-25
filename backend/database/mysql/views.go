package mysql

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/NeverStopDreamingWang/goi"
	"github.com/NeverStopDreamingWang/goi/db"
	"github.com/NeverStopDreamingWang/goi/db/mysql"
)

// 参数验证
type listValidParams struct {
	Page         int             `name:"page" type:"int" required:"true"`
	PageSize     int             `name:"page_size" type:"int" required:"true"`
	Username     *string         `name:"username" type:"string"`
	Status       *UserStatusType `name:"status" type:"int"`
	RoleId       *int64          `name:"role_id" type:"int"`
	DepartmentId *int64          `name:"department_id" type:"int"`
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

	mysqlDB := db.Connect[*mysql.Engine]("mysql")

	mysqlDB.SetModel(UserModel{}) // 设置操作表

	if params.Username != nil {
		mysqlDB = mysqlDB.Where("`username` like ?", "%"+*params.Username+"%")
	}
	if params.Status != nil {
		mysqlDB = mysqlDB.Where("`status` = ?", params.Status)
	}
	if params.RoleId != nil {
		mysqlDB = mysqlDB.Where("`role_id` = ?", params.RoleId)
	}
	if params.DepartmentId != nil {
		mysqlDB = mysqlDB.Where("`department_id` = ?", params.DepartmentId)
	}
	total, total_page, err := mysqlDB.Page(params.Page, params.PageSize)
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询用户失败",
			Results: err.Error(),
		}
	}

	var user_list []*UserModel
	err = mysqlDB.Select(&user_list)
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询用户失败",
			Results: err.Error(),
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

	mysqlDB := db.Connect[*mysql.Engine]("mysql")

	user := &UserModel{}
	mysqlDB.SetModel(UserModel{})
	err := mysqlDB.Where("`id` = ?", pk).First(user)
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

	mysqlDB := db.Connect[*mysql.Engine]("mysql")

	instance := &UserModel{}
	mysqlDB.SetModel(UserModel{})
	err := mysqlDB.Where("`id` = ?", pk).First(instance)
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

	mysqlDB := db.Connect[*mysql.Engine]("mysql")

	instance := &UserModel{}
	mysqlDB.SetModel(UserModel{})
	err := mysqlDB.Where("`id` = ?", pk).First(instance)
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
