package basic

import (
	"goi_example/backend/goi_example"
	"goi_example/server/web"

	"github.com/NeverStopDreamingWang/goi"
)

func init() {
	// 子路由
	basicRouter := goi_example.Server.Router.Include("basic/", "父路由")
	{
		// 注册一个路径
		basicRouter.Path("test1", "测试路由1", goi.ViewSet{GET: Test1})

		// 创建一个三级子路由
		paramsRouter := basicRouter.Include("params/", "传参")
		{
			paramsRouter.Path("path/<string:name>", "路由传参", goi.ViewSet{GET: TestPathParams})

			paramsRouter.Path("query", "query 查询字符串传参", goi.ViewSet{GET: TestQueryParams})

			paramsRouter.Path("body", "body 请求体传参", goi.ViewSet{GET: TestBodyParams})

			paramsRouter.Path("valid", "参数验证", goi.ViewSet{POST: TestParamsValid})

			paramsRouter.Path("converter/<string:name>", "自定义路由转换器获取参数", goi.ViewSet{GET: TestConverterParamsStr})

			paramsRouter.Path("converter_phone/<phone:phone>", "测试 phone 路由转换器", goi.ViewSet{GET: TestConverterParamsPhone})

			paramsRouter.Path("context/<string:name>", "从上下文中获取数据", goi.ViewSet{GET: TestContext})

			paramsRouter.Path("context_params", "从上下文参数中获取数据", goi.ViewSet{GET: TestContextParams})
		}
		// 静态路由
		// 静态文件
		basicRouter.StaticFile("test.txt", "测试静态文件", "static/test.txt")

		// 静态目录
		basicRouter.StaticDir("static/", "测试静态目录", "static")

		// embed.FS 静态文件
		basicRouter.StaticFileFS("test1.txt", "测试静态FS文件", web.Assets, "assets/test.txt")

		// embed.FS 静态目录
		basicRouter.StaticDirFS("assets/", "测试静态FS目录", web.Assets, "assets")

		// 自定义方法
		// 返回文件
		basicRouter.Path("test_file", "返回文件", goi.ViewSet{GET: TestFile})

		// 异常处理
		basicRouter.Path("test_panic", "异常处理", goi.ViewSet{GET: TestPanic})

		// 缓存
		basicRouter.Path("cache", "测试缓存", goi.ViewSet{GET: TestCacheGet, POST: TestCacheSet, DELETE: TestCacheDel})
	}
}
