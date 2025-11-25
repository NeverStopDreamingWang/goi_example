package basic

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"goi_example/backend/permission/user"

	"github.com/NeverStopDreamingWang/goi"
)

func Test1(request *goi.Request) any {
	goi.Log.DebugF("Test1")

	return map[string]any{
		"status": http.StatusOK,
		"msg":    "Hello World",
		"data":   "",
	}
}

// 路由传参
func TestPathParams(request *goi.Request) any {
	var name string
	var validationErr goi.ValidationError
	validationErr = request.PathParams.Get("name", &name) // 路由传参
	if validationErr != nil {
		return validationErr.Response()
	}
	goi.Log.DebugF("TestPathParams 参数: %v 参数类型: %T", name, name)

	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: name,
	}
}

type testQueryParamsValidParams struct {
	Name string `name:"name" type:"string" required:"true"`
	Age  *int64 `name:"age" type:"int"`
}

// query 查询字符串传参
func TestQueryParams(request *goi.Request) any {
	var params testQueryParamsValidParams
	var queryParams goi.Params
	var validationErr goi.ValidationError
	queryParams = request.QueryParams()

	// 获取一个参数
	var name string
	validationErr = queryParams.Get("name", &name)
	if validationErr != nil {
		return validationErr.Response()
	}
	goi.Log.DebugF("TestQueryParams 参数: %v 参数类型: %T", name, name)

	// 获取多个参数
	validationErr = queryParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}
	goi.Log.DebugF("TestQueryParams 参数: %+v 参数类型: %T", params, params)

	return goi.Response{
		Status: http.StatusOK, // 返回响应状态码
		Data: goi.Data{ // 响应数据
			Code:    http.StatusOK, // 返回自定义状态码
			Message: "",
			Results: params,
		},
	}
}

type testBodyParamsValidParams struct {
	Name *string `name:"name" type:"string" required:"true"`
	Age  *int64  `name:"age" type:"int"`
}

// body 请求体传参
func TestBodyParams(request *goi.Request) any {
	var bodyParams goi.Params
	var validationErr goi.ValidationError
	bodyParams = request.BodyParams()

	// 获取一个参数
	var name string
	validationErr = bodyParams.Get("name", &name)
	if validationErr != nil {
		return validationErr.Response()
	}
	goi.Log.DebugF("TestBodyParams 参数: %v 参数类型: %T", name, name)

	// 获取多个参数
	var params testBodyParamsValidParams
	validationErr = bodyParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}
	goi.Log.DebugF("TestBodyParams 参数: %+v 参数类型: %T", params, params)

	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: params,
	}
}

// 参数验证
type testParamsValidParams struct {
	Username string            `name:"username" type:"string" required:"true"`
	Password string            `name:"password" type:"string" required:"true"`
	Age      int64             `name:"age" type:"int" required:"true"`
	Phone    string            `name:"phone" type:"phone" required:"true"`
	Args     []int64           `name:"args" type:"slice"`
	Kwargs   map[string]string `name:"kwargs" type:"map"`
}

// type:"name" 字段类型, name 验证器名称
// required:"bool" 字段是否必填，bool 布尔值： false 默认可选，true 必传参数
// allow_null:"bool" 字段值是否允许为空，bool 布尔值，当 required == true 时默认 allow_null=false 不允许为空，否则默认 allow_null=true 允许为空
// 支持
// int *int []*int []... map[string]*int map[...]...
// ...

func TestParamsValid(request *goi.Request) any {
	var params testParamsValidParams
	var bodyParams goi.Params
	var validationErr goi.ValidationError
	bodyParams = request.BodyParams() // Body 传参
	validationErr = bodyParams.ParseParams(&params)
	if validationErr != nil {
		// 验证器返回
		return validationErr.Response()

		// 自定义返回
		// return goi.Response{
		// 	Status: http.StatusOK,
		// 	Data: goi.Data{
		// 		Code: http.StatusBadRequest,
		// 		Message:    "参数错误",
		// 		Results:   nil,
		// 	},
		// }
	}
	goi.Log.Debug("TestParamsValid", params)

	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: params,
	}
}

func TestConverterParamsStr(request *goi.Request) any {
	var name string
	var validationErr goi.ValidationError
	validationErr = request.PathParams.Get("name", &name)
	if validationErr != nil {
		return validationErr.Response()
	}
	goi.Log.DebugF("TestConverterParamsStr 参数: %v 参数类型:  %T", name, name)

	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: name,
	}
}

// 测试手机号路由转换器
func TestConverterParamsPhone(request *goi.Request) any {
	var phone string
	var validationErr goi.ValidationError
	validationErr = request.PathParams.Get("phone", &phone)
	if validationErr != nil {
		return validationErr.Response()
	}
	goi.Log.DebugF("TestConverterParamsPhone 参数: %v 参数类型:  %T", phone, phone)

	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: phone,
	}
}

// 上下文
func TestContext(request *goi.Request) any {
	// 请求上下文
	requestID := request.Object.Context().Value("requestID")

	goi.Log.Debug("requestID", requestID)

	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: requestID,
	}
}

// 上下文参数
func TestContextParams(request *goi.Request) any {
	var validationErr goi.ValidationError

	// 请求上下文参数
	userInfo := &user.UserModel{}
	validationErr = request.Params.Get("user", userInfo)
	if validationErr != nil {
		return validationErr.Response()
	}
	goi.Log.Debug("user", userInfo)

	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: userInfo,
	}
}

// 返回文件
func TestFile(request *goi.Request) any {
	absolutePath := filepath.Join(goi.Settings.BASE_DIR, "static/test.txt")
	file, err := os.Open(absolutePath)
	if err != nil {
		_ = file.Close()
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "读取文件失败",
		}
	}
	// return file // 返回文件对象
	return file // 返回文件对象
}

// 异常处理
func TestPanic(request *goi.Request) any {
	var bodyParams goi.Params
	bodyParams = request.BodyParams()
	name := bodyParams["name"]

	msg := fmt.Sprintf("参数: %v 参数类型:  %T", name, name)
	goi.Log.Debug("TestPanic", msg)

	panic(name)

	return goi.Data{
		Code:    http.StatusOK,
		Message: msg,
		Results: "",
	}
}
