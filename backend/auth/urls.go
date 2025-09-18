package auth

import (
	"goi_example/backend/goi_example"

	"github.com/NeverStopDreamingWang/goi"
)

func init() {
	// 认证
	authRouter := goi_example.Server.Router.Include("auth/", "认证")
	{
		authRouter.Path("captcha", "获取图片验证码", goi.ViewSet{GET: captchaView})
		authRouter.Path("login", "用户登录", goi.ViewSet{POST: loginView})
	}
}
