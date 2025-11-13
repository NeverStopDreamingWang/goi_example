package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"goi_example/backend/goi_example"
	"goi_example/backend/permission/user"
	"goi_example/backend/utils"

	"github.com/NeverStopDreamingWang/goi"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	goi_example.Server.MiddleWare = append(goi_example.Server.MiddleWare, &AuthMiddleWare{})
}

// Token
type authValidParams struct {
	Token *string `name:"token" type:"string"`
}

type AuthMiddleWare struct{}

func (AuthMiddleWare) ProcessRequest(request *goi.Request) interface{} {
	// fmt.Println("请求中间件", request.Object.URL)

	// 跳过验证
	if strings.HasPrefix(request.Object.URL.Path, "/api") == false &&
		strings.HasPrefix(request.Object.URL.Path, goi_example.STATIC_URL) == false {
		return nil
	}

	var apiList = []string{
		"/api/auth",  // 认证
		"/api/basic", // 基础
	}
	for _, api := range apiList {
		// 跳过验证
		if strings.HasPrefix(request.Object.URL.Path, api) {
			return nil
		}
	}

	token := request.Object.Header.Get("Authorization")
	if token == "" {
		var params authValidParams
		var queryParams goi.Params
		var validationErr goi.ValidationError
		queryParams = request.QueryParams() // Query 传参
		validationErr = queryParams.ParseParams(&params)
		if validationErr != nil {
			return validationErr.Response()
		}
		if params.Token != nil {
			token = *params.Token
		}
	}

	payloads := &utils.Payloads{}
	err := utils.CheckToken(token, goi.Settings.SECRET_KEY, payloads)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return goi.Data{
				Code:    http.StatusUnauthorized,
				Message: "token 解码错误",
				Results: err,
			}
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return goi.Data{
				Code:    http.StatusUnauthorized,
				Message: "token 已过期",
				Results: err,
			}
		}
		return goi.Data{
			Code:    http.StatusUnauthorized,
			Message: "token 验证失败",
			Results: err,
		}
	}
	userInfo, err := user.GetUser(payloads.User_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return goi.Data{
				Code:    http.StatusUnauthorized,
				Message: "用户不存在",
				Results: nil,
			}
		}
		return goi.Data{
			Code:    http.StatusUnauthorized,
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

	// 写入请求参数
	request.Params.Set("user", userInfo)
	return nil
}

func (AuthMiddleWare) ProcessView(request *goi.Request) interface{} { return nil }

func (AuthMiddleWare) ProcessException(request *goi.Request, exception any) interface{} { return nil }

func (AuthMiddleWare) ProcessResponse(request *goi.Request, response *goi.Response) {}
