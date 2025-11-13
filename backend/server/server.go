package server

import (
	"fmt"
	"reflect"
	"strings"

	"goi_example/backend/goi_example"

	"github.com/NeverStopDreamingWang/goi"

	// 注册app
	_ "goi_example/backend/auth"             // 认证
	_ "goi_example/backend/basic"            // 基础
	_ "goi_example/backend/database"         // 数据库
	_ "goi_example/backend/database/mongodb" // mongodb
	_ "goi_example/backend/database/mysql"   // mysql
	_ "goi_example/backend/permission"       // 权限
	_ "goi_example/backend/permission/role"  // 角色
	_ "goi_example/backend/permission/user"  // 用户

	// 注册中间件
	_ "goi_example/backend/middleware"
)

func Start() {
	// 获取所有路径信息
	fmt.Println("Route:")
	route := goi_example.Server.Router.GetRoute()
	printRouteInfo(route, 1)

	// 启动服务
	goi_example.Server.RunServer()
}

func Stop() {
	err := goi_example.Server.StopServer()
	panic(err)
}

var replaceStr = "  "

func printRouteInfo(route goi.Route, count int) {
	fmt.Printf("%vPath: %v Desc: %v Pattern: %v Regex:%v ParamInfos:%+v\n",
		strings.Repeat(replaceStr, count),
		route.Path,
		route.Desc,
		route.Pattern,
		route.Regex,
		route.ParamInfos,
	)
	var method_list []string
	// ViewSet
	ViewSetValue := reflect.ValueOf(route.ViewSet)
	ViewSetType := ViewSetValue.Type()
	// 获取字段
	for i := 0; i < ViewSetType.NumField(); i++ {
		method := ViewSetType.Field(i).Name
		if ViewSetValue.Field(i).IsZero() {
			continue
		}
		method_list = append(method_list, method)
	}
	fmt.Printf("%vViewSet: %v\n", strings.Repeat(replaceStr, count+1), strings.Join(method_list, ", "))
	if len(route.Children) > 0 {
		fmt.Printf("%vChildren:\n", strings.Repeat(replaceStr, count))
		for _, child := range route.Children {
			printRouteInfo(child, count+1)
		}
	}
}
