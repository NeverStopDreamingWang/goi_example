package user

import (
	"database/sql"
	"errors"
	"net/http"

	"goi_example/backend/permission/role"

	"github.com/NeverStopDreamingWang/goi"
	"github.com/NeverStopDreamingWang/goi/auth"
	"github.com/NeverStopDreamingWang/goi/db"
	"github.com/NeverStopDreamingWang/goi/db/sqlite3"
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
	Menu_List     []*userMenuList `json:"menu_list"`
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

	user.Role = &role.RoleModel{}
	sqlite3DB.SetModel(role.RoleModel{})
	err = sqlite3DB.Where("`id` = ?", user.RoleId).First(user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return goi.Data{
				Code:    http.StatusBadRequest,
				Message: "角色不存在",
				Results: nil,
			}
		}
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询数据库错误",
			Results: err.Error(),
		}
	}

	roleMenuList := []role.RoleMenuModel{}
	sqlite3DB.SetModel(role.RoleMenuModel{})
	err = sqlite3DB.Where("role_id = ?", user.RoleId).Select(&roleMenuList)
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询数据库错误",
			Results: err.Error(),
		}
	}

	MenuList := make([]*userMenuList, len(roleMenuList))

	sqlite3DB.SetModel(role.MenuModel{})
	for i, roleMenu := range roleMenuList {
		MenuList[i] = &userMenuList{}
		err = sqlite3DB.Where("id = ?", roleMenu.Menu_Id).First(MenuList[i])
		if err != nil {
			return goi.Data{
				Code:    http.StatusInternalServerError,
				Message: "查询数据库错误",
				Results: err.Error(),
			}
		}
	}

	user.Menu_List = get_children_menu(MenuList)
	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: user,
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
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "用户不存在",
			Results: err.Error(),
		}
	}
	if params.OldPassword != nil {
		if auth.CheckPassword(*params.OldPassword, *instance.Password) == false {
			return goi.Data{
				Code:    http.StatusBadRequest,
				Message: "旧密码错误",
				Results: nil,
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
