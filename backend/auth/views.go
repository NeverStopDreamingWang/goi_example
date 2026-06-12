package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/NeverStopDreamingWang/goi/v2"
	"github.com/NeverStopDreamingWang/goi/v2/db/sqlite3"
	"github.com/NeverStopDreamingWang/goi/v2/response"

	"goi_example/backend/permission/user"
	"goi_example/backend/utils"
	"goi_example/backend/utils/captcha"

	"github.com/NeverStopDreamingWang/goi/v2/auth"
	"github.com/NeverStopDreamingWang/goi/v2/db"
)

func captchaView(request *goi.Request) any {
	id, base64Str, err := captcha.NewCaptcha()
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "获取验证码错误",
			Data:    nil,
		}
	}
	return response.Data{
		Code:    http.StatusOK,
		Message: "",
		Data: map[string]string{
			"id":    id,
			"image": base64Str,
		},
	}
}

type loginParams struct {
	Email     *string `name:"email" type:"string"`
	Username  *string `name:"username" type:"string"`
	Password  string  `name:"password" type:"string" required:"true"`
	CaptchaId string  `name:"captcha_id" type:"string"`
	Captcha   string  `name:"captcha" type:"string"`
}

func loginView(request *goi.Request) any {
	var params loginParams
	var bodyParams goi.Params
	var validationErr goi.ValidationError
	var err error

	bodyParams = request.BodyParams()
	validationErr = bodyParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	if params.CaptchaId != "" {
		if params.Captcha == "" {
			return response.Data{
				Code:    http.StatusBadRequest,
				Message: "请输入验证码",
				Data:    nil,
			}
		}
		err = captcha.VerifyCode(params.CaptchaId, params.Captcha)
		if err != nil {
			return response.Data{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Data:    nil,
			}
		}

	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	userObject := user.UserModel{}
	sqlite3DB.SetModel(user.UserModel{})
	if params.Username != nil && *params.Username != "" {
		err = sqlite3DB.Where("username=?", params.Username).First(&userObject)
	} else if params.Email != nil && *params.Email != "" {
		err = sqlite3DB.Where("email=?", params.Email).First(&userObject)
	} else {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: "请输入账号或邮箱",
			Data:    nil,
		}
	}

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

	if *userObject.Status == user.DISABLE {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: "当前账号已被禁用",
			Data:    nil,
		}
	}

	if auth.CheckPassword(params.Password, *userObject.Password) == false {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: "账号或密码错误",
			Data:    nil,
		}
	}

	payload := utils.Payloads{
		Exp:      time.Now().In(goi.GetLocation()).Add(2 * time.Hour).Unix(), // 设置过期时间为2小时后
		UserId:   *userObject.Id,
		Username: *userObject.Username,
	}
	token, err := utils.NewToken(payload, goi.Settings.SecretKey)
	if err != nil {
		goi.Log.Error("生成 Token 错误", err.Error())
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "生成 Token 错误",
			Data:    err,
		}
	}

	// 更新最后登录时间
	lastLoginTime := goi.GetTime()
	err = userObject.Update(&user.UserModel{
		LastLoginTime: &lastLoginTime,
	})
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	data := map[string]any{
		"user":  userObject,
		"token": token,
	}
	return response.Data{
		Code:    http.StatusOK,
		Message: "登录成功",
		Data:    data,
	}
}
