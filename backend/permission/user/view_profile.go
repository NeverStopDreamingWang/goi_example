package user

import (
	"database/sql"
	"errors"
	"net/http"

	"goi_example/backend/permission/role"

	"github.com/NeverStopDreamingWang/goi/v2"
	"github.com/NeverStopDreamingWang/goi/v2/auth"
	"github.com/NeverStopDreamingWang/goi/v2/db"
	"github.com/NeverStopDreamingWang/goi/v2/db/sqlite3"
	"github.com/NeverStopDreamingWang/goi/v2/response"
)

type userMenuList struct {
	Id       *int64          `json:"id"`
	ParentId *int64          `json:"parent_id"`
	Name     *string         `json:"label"`
	Icon     *string         `json:"icon"`
	Path     *string         `json:"index"`
	Children []*userMenuList `json:"children"`
}
type profileUser struct {
	Id            *int64          `json:"id"`
	Username      *string         `json:"username"`
	Email         *string         `json:"email"`
	Status        *UserStatusType `json:"status"`
	RoleId        *int64          `json:"role_id"`
	LastLoginTime *string         `json:"last_login_time"`
	CreateTime    *string         `json:"create_time"`
	UpdateTime    *string         `json:"update_time"`
	Role          *role.RoleModel `json:"role"`
	MenuList      []*userMenuList `json:"menu_list"`
}

func profileRetrieveView(request *goi.Request) any {
	// 获取当前用户信息
	var userObject UserModel
	validationErr := request.Params.Get("user", &userObject)
	if validationErr != nil {
		return validationErr.Response()
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")
	user := profileUser{}
	sqlite3DB.SetModel(UserModel{})
	err := sqlite3DB.Where("`id` = ?", userObject.Id).First(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.Data{
				Code:    http.StatusBadRequest,
				Message: "用户不存在",
				Data:    nil,
			}
		}
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询数据库错误",
			Data:    err.Error(),
		}
	}

	user.Role = &role.RoleModel{}
	sqlite3DB.SetModel(role.RoleModel{})
	err = sqlite3DB.Where("`id` = ?", user.RoleId).First(user.Role)
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

	roleMenuList := []role.RoleMenuModel{}
	sqlite3DB.SetModel(role.RoleMenuModel{})
	err = sqlite3DB.Where("role_id = ?", user.RoleId).Select(&roleMenuList)
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询数据库错误",
			Data:    err.Error(),
		}
	}

	MenuList := make([]*userMenuList, len(roleMenuList))

	sqlite3DB.SetModel(role.MenuModel{})
	for i, roleMenu := range roleMenuList {
		MenuList[i] = &userMenuList{}
		err = sqlite3DB.Where("id = ?", roleMenu.MenuId).First(MenuList[i])
		if err != nil {
			return response.Data{
				Code:    http.StatusInternalServerError,
				Message: "查询数据库错误",
				Data:    err.Error(),
			}
		}
	}

	user.MenuList = get_children_menu(MenuList)
	return response.Data{
		Code:    http.StatusOK,
		Message: "",
		Data:    user,
	}
}

// 参数验证
type profileUpdateValidParams struct {
	Username    *string `name:"username" type:"string"`
	OldPassword *string `name:"old_password" type:"string"`
	NewPassword *string `name:"new_password" type:"string"`
	Email       *string `name:"email" type:"string"`
}

func profileUpdateView(request *goi.Request) any {
	var params profileUpdateValidParams
	var bodyParams goi.Params
	var validationErr goi.ValidationError

	bodyParams = request.BodyParams() // Body 传参
	validationErr = bodyParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	// 获取当前用户信息
	var userObject UserModel
	validationErr = request.Params.Get("user", &userObject)
	if validationErr != nil {
		return validationErr.Response()
	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	instance := &UserModel{}
	sqlite3DB.SetModel(UserModel{})
	err := sqlite3DB.Where("`id` = ?", userObject.Id).First(instance)
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "用户不存在",
			Data:    err.Error(),
		}
	}
	if params.OldPassword != nil {
		if auth.CheckPassword(*params.OldPassword, *instance.Password) == false {
			return response.Data{
				Code:    http.StatusBadRequest,
				Message: "旧密码错误",
				Data:    nil,
			}
		}
	}

	user := &UserModel{
		Id:       instance.Id,
		Username: params.Username,
		Password: params.NewPassword,
		Email:    params.Email,
	}

	// 参数验证
	err = user.Validate()
	if err != nil {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}
	// 更新
	err = instance.Update(user)
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
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
