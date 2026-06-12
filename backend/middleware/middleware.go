package middleware

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"goi_example/backend/goi_example"
	"goi_example/backend/permission/user"
	"goi_example/backend/utils"

	"github.com/NeverStopDreamingWang/goi/v2"
	"github.com/NeverStopDreamingWang/goi/v2/response"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	goi_example.ApiRouter.Use(&AuthMiddleware{})
}

// Token
type authValidParams struct {
	Token *string `name:"token" type:"string"`
}

type AuthMiddleware struct{}

func (AuthMiddleware) ProcessRequest(request *goi.Request) any {
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
	err := utils.CheckToken(token, goi.Settings.SecretKey, payloads)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return response.Data{
				Code:    http.StatusUnauthorized,
				Message: "token 解码错误",
				Data:    err,
			}
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return response.Data{
				Code:    http.StatusUnauthorized,
				Message: "token 已过期",
				Data:    err,
			}
		}
		return response.Data{
			Code:    http.StatusUnauthorized,
			Message: "token 验证失败",
			Data:    err,
		}
	}
	userInfo, err := user.GetUser(payloads.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response.Data{
				Code:    http.StatusUnauthorized,
				Message: "用户不存在",
				Data:    nil,
			}
		}
		return response.Data{
			Code:    http.StatusUnauthorized,
			Message: "查询数据库错误",
			Data:    err.Error(),
		}
	}

	if *userInfo.Status == user.DISABLE {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: "当前账号已被禁用",
			Data:    nil,
		}
	}

	// 写入请求参数
	request.Params.Set("user", userInfo)
	return nil
}

func (AuthMiddleware) ProcessView(request *goi.Request) any { return nil }

func (AuthMiddleware) ProcessException(request *goi.Request, exception any) any { return nil }

func (AuthMiddleware) ProcessResponse(request *goi.Request, response *goi.Response) {}
