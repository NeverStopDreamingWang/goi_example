package auth

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/NeverStopDreamingWang/goi"
	"github.com/NeverStopDreamingWang/goi/db/sqlite3"

	"goi_example/backend/permission/user"
	"goi_example/backend/utils"
	"goi_example/backend/utils/captcha"

	"github.com/NeverStopDreamingWang/goi/auth"
	"github.com/NeverStopDreamingWang/goi/db"
	"github.com/NeverStopDreamingWang/goi/jwt"
)

func captchaView(request *goi.Request) interface{} {
	id, base64Str, err := captcha.NewCaptcha()
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "获取验证码错误",
			Results: nil,
		}
	}
	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: map[string]string{
			"id":    id,
			"image": base64Str,
		},
	}
}

type loginParams struct {
	Email      *string `name:"email" type:"string"`
	Username   *string `name:"username" type:"string"`
	Password   string  `name:"password" type:"string" required:"true"`
	Captcha_id string  `name:"captcha_id" type:"string"`
	Captcha    string  `name:"captcha" type:"string"`
}

func loginView(request *goi.Request) interface{} {
	var params loginParams
	var bodyParams goi.Params
	var validationErr goi.ValidationError
	var err error

	bodyParams = request.BodyParams()
	validationErr = bodyParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	if params.Captcha_id != "" {
		if params.Captcha == "" {
			return goi.Data{
				Code:    http.StatusBadRequest,
				Message: "请输入验证码",
				Results: nil,
			}
		}
		err = captcha.VerifyCode(params.Captcha_id, params.Captcha)
		if err != nil {
			return goi.Data{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Results: nil,
			}
		}

	}

	sqlite3DB := db.Connect[*sqlite3.Engine]("default")

	userInfo := user.UserModel{}
	sqlite3DB.SetModel(user.UserModel{})
	if params.Username != nil && *params.Username != "" {
		err = sqlite3DB.Where("username=?", params.Username).First(&userInfo)
	} else if params.Email != nil && *params.Email != "" {
		err = sqlite3DB.Where("email=?", params.Email).First(&userInfo)
	} else {
		return goi.Data{
			Code:    http.StatusBadRequest,
			Message: "请输入账号或邮箱",
			Results: nil,
		}
	}

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

	if *userInfo.Status == user.DISABLE {
		return goi.Data{
			Code:    http.StatusBadRequest,
			Message: "当前账号已被禁用",
			Results: nil,
		}
	}

	if auth.CheckPassword(params.Password, *userInfo.Password) == false {
		return goi.Data{
			Code:    http.StatusBadRequest,
			Message: "账号或密码错误",
			Results: nil,
		}
	}

	header := jwt.Header{
		Alg: jwt.AlgHS256,
		Typ: jwt.TypJWT,
	}

	// 设置过期时间
	twoHoursLater := goi.GetTime().Add(24 * 15 * time.Hour)

	payloads := utils.Payloads{ // 包含 jwt.Payloads
		Payloads: jwt.Payloads{
			Exp: jwt.ExpTime{twoHoursLater},
		},
		User_id:  *userInfo.Id,
		Username: *userInfo.Username,
	}
	token, err := jwt.NewJWT(header, payloads, goi.Settings.SECRET_KEY)
	if err != nil {
		goi.Log.Error("生成 Token 错误", err.Error())
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "生成 Token 错误",
			Results: err,
		}
	}

	data := map[string]interface{}{
		"user":  userInfo,
		"token": token,
	}
	return goi.Data{
		Code:    http.StatusOK,
		Message: "登录成功",
		Results: data,
	}
}
