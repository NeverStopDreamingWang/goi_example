# Goi Example

一个用于存放 Goi 示例和教程的资源库。

当前示例依赖 **Goi v3.0.3**（`github.com/NeverStopDreamingWang/goi/v3`）。下面的导入和 ORM 调用均按 v3 接口编写；MongoDB 驱动仍使用其自身的 `go.mongodb.org/mongo-driver/v2`，不要一起替换。

## 项目结构

```
goi_example/
├── backend/           # 后端代码
│   ├── server/        # 服务管理：启动、停止
│   ├── goi_example/   # 项目配置
│   ├── auth/          # 认证模块
│   ├── basic/         # 基础功能
│   ├── middleware/    # 中间件
│   ├── database/      # 数据库相关
│   │   ├── urls.go    # 路由组
│   │   ├── mongodb/   # mongodb 数据库相关
│   │   └── mysql/     # mysql 数据库相关
│   └── permission/    # 权限管理
│       ├── urls.go    # 路由组
│       ├── role/      # 角色管理
│       └── user/      # 用户管理
├── server/            # 服务入口
│   ├── cmd/           # 命令行入口
│   ├── web/           # 前端资源
│   └── main.go        # 主入口
├── build/             # 构建输出目录
│   ├── config.yaml    # 配置文件
│   ├── install.sh     # 安装脚本
│   └── server_test.go # 服务测试
├── Makefile           # 构建脚本
└── README.md          # 项目说明
```

## 快速开始

### 前置条件

需要 Go 1.24 或更高版本。默认 SQLite 使用 `github.com/mattn/go-sqlite3`，需启用 CGO 并安装可用的 C 编译器。

### 运行本示例项目

在仓库根目录运行。首次使用时，将 `build/config.yaml` 复制为根目录的 `config.yaml`，按需修改端口和数据库配置；已有配置不要覆盖。

Linux/macOS：

```sh
cp build/config.yaml config.yaml
export CGO_ENABLED=1
go mod download
go run ./server
```

Windows PowerShell：

```powershell
Copy-Item build/config.yaml config.yaml
$env:CGO_ENABLED = "1"
go mod download
go run ./server
```

实际入口为 `server/main.go`，默认命令会调用 `backend/server.Start()`，不需要额外的 `web` 或 `start` 子命令。配置从当前工作目录的 `config.yaml` 读取，所以应在仓库根目录启动。默认端口为 8080，可访问 `http://127.0.0.1:8080/`。

默认配置只启用 SQLite；启用 MySQL、Redis 或 MongoDB 时，需要自行准备服务和凭据，并检查 MySQL 模型的 `Settings.Schema` 与目标数据库是否一致。迁移不会自动创建数据库。本仓库目前仍在包 `init` 中加载配置、注册路由并迁移 SQL 模型，因此导入这些包就可能访问数据库；不要把普通包测试或帮助命令视为完全无副作用的检查。示例密钥与初始化账号仅供本地演示，不应直接用于生产。

### 本地验证（不启动服务）

在仓库根目录运行以下已验证的检查命令。SQLite 测试仍需启用 CGO，并安装可用的 C 编译器：

```sh
go build ./...
go vet ./...
go test ./backend/utils/...
go test ./backend/goi_example/validator.go ./backend/goi_example/validator_test.go
```

最后一条使用显式文件列表，仅测试验证器，避开项目配置和应用的 `init` 初始化；不要直接替换为整个配置包的测试。普通 `go test ./...` 包含会启动服务器的测试，不作为这里的无监听检查命令。

### 获取 Goi

通过 Go 的模块支持，当你在代码中添加导入时，`go [build|run|test]` 会自动获取必要的依赖：

```go
import "github.com/NeverStopDreamingWang/goi/v3"
```

或者，使用 `go get`：

```bash
go get github.com/NeverStopDreamingWang/goi/v3@v3.0.3
```

### 独立最小示例

一个基础示例：

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/NeverStopDreamingWang/goi/v3"
	"github.com/NeverStopDreamingWang/goi/v3/response"
)

func Ping(request *goi.Request) any {
	goi.Log.DebugF("Test1")

	return response.Data{
		Code:    http.StatusOK,
		Message: "Hello World",
		Data:    nil,
	}
}

func main() {
	// 创建 Goi 服务器
	Server := goi.NewHTTPServer()
	// 创建并注册日志切割任务
	Server.Log = goi.NewLogger("server.log")
	goi.Log = Server.Log
	goi.RegisterOnStartup(Server.Log)
	// RunServer 完成全部停机流程后，由使用方关闭日志文件
	defer func() {
		if err := Server.Log.Close(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}()

	// 网络协议
	Server.Settings.Network = "tcp" // 默认 "tcp" 常用网络协议 "tcp"、"tcp4"、"tcp6"、"udp"、"udp4"、"udp6"
	// 监听地址 0.0.0.0
	Server.Settings.BindAddress = "0.0.0.0" // 默认 127.0.0.1
	// 监听端口 8080
	Server.Settings.Port = 8080

	// 注册路由
	Server.Router.Path("ping", "测试接口", goi.ViewSet{GET: Ping})

	// 启动服务器
	Server.RunServer()
}
```

下面命令用于上面的独立 `main.go` 示例；运行本仓库完整项目仍使用 `go run ./server`：

```bash
go run main.go
```

然后在浏览器中访问 `http://0.0.0.0:8080/ping` 查看响应！

## 路由

### 多级路由

子父路由，内置路由转换器，自定义路由转换器

* `Router.Include(path string, desc string)` 创建一个子路由
    * `path` 路由
    * `desc` 描述

* `Router.Path(path string, desc string, view ViewSet)` 注册一个路由
    * `path` 路由
    * `desc` 描述
    * `view` 视图 `ViewSet` 类型
        * `GET`      `HandlerFunc` 获取资源
        * `HEAD`     `HandlerFunc` 获取响应头
        * `POST`     `HandlerFunc` 创建资源
        * `PUT`      `HandlerFunc` 更新资源
        * `PATCH`    `HandlerFunc` 部分更新
        * `DELETE`   `HandlerFunc` 删除资源
        * `CONNECT`  `HandlerFunc` 建立隧道
        * `OPTIONS`  `HandlerFunc` 查询支持的方法
        * `TRACE`    `HandlerFunc` 回显请求
        * `NoMethod` `HandlerFunc` 未注册或不支持的HTTP方法时的处理函数，默认返回 405

### 静态路由

静态路由映射，返回文件对象即响应该文件内容

静态路由映射可通过两种方式注册：

方式一：Router.* 方法（内部基于公共 ViewSet 生成器）

* `Router.StaticFile(path string, desc string, filePath string)` 注册静态文件路由
    * `path` 路由
    * `desc` 描述
    * `filePath` 文件路径，支持`相对路径`和`绝对路径`

* `Router.StaticDir(path string, desc string, dirPath http.Dir)` 注册静态目录路由
    * `path` 路由
    * `desc` 描述
    * `dirPath` 静态映射路径，支持`相对路径`和`绝对路径`

* `Router.StaticFileFS(path string, desc string, fileFS embed.FS, defaultPath string)` 注册 `embed.FS` 静态文件路由
    * `path` 路由
    * `desc` 描述
    * `fileFS` 嵌入式文件系统，不是磁盘路径
    * `defaultPath` 文件系统内的相对文件路径


* `Router.StaticDirFS(path string, desc string, dirFS embed.FS, basePath string)` 注册 `embed.FS` 静态目录路由
    * `path` 路由
    * `desc` 描述
    * `dirFS` 嵌入式文件系统，不是磁盘路径
    * `basePath` 文件系统内的相对根目录，空字符串按 `.` 处理
    * 说明:
        - `path` 路径之后自动添加 `/<path:fileName>` 参数
        - 注册时通过 `fs.Sub` 将 `basePath` 锚定为文件系统根目录，再按 `fileName` 查找文件

### 静态 View 生成器

提供四个公共方法生成静态资源的 HandlerFunc，配合 `Router.Path` 使用：

- `goi.StaticFileView(filePath string) HandlerFunc`
- `goi.StaticDirView(dirPath http.Dir) HandlerFunc`
- `goi.StaticFileFSView(fileFS fs.FS, defaultPath string) HandlerFunc`
- `goi.StaticDirFSView(dirFS fs.FS) HandlerFunc`

目录生成器接收已锚定的 `fs.FS`，不再单独接收基础路径。下面示例中的 `fs.Sub` 需要导入 `io/fs`；`router`、`assets` 分别表示已创建的路由和嵌入资源包。

```go
// 单文件
router.Path("test.txt", "测试文件", goi.ViewSet{GET: goi.StaticFileView("static/test.txt")})

// 目录（带路由参数 <path:fileName>）
router.Path("static/<path:fileName>", "静态目录", goi.ViewSet{GET: goi.StaticDirView(http.Dir("static"))})

// 嵌入式单文件
router.Path("logo.svg", "静态FS文件", goi.ViewSet{GET: goi.StaticFileFSView(assets.LogoFS, "logo.svg")})
// 嵌入式目录（带路由参数 <path:fileName>）
assetsRoot, err := fs.Sub(assets.AssetsFS, "assets")
if err != nil {
	panic(err)
}
router.Path("assets/<path:fileName>", "静态FS目录", goi.ViewSet{GET: goi.StaticDirFSView(assetsRoot)})
```

## 中间件

框架使用 `goi.Middleware` 接口，中间件必须实现全部三个方法；不需要处理的阶段可提供空实现：

```go
package middleware

import (
	"github.com/NeverStopDreamingWang/goi/v3"
)

type MyMiddleware struct{}

var _ goi.Middleware = (*MyMiddleware)(nil)

func (*MyMiddleware) ProcessRequest(request *goi.Request) any {
	return nil
}

func (*MyMiddleware) ProcessException(request *goi.Request, exception any) any {
	return nil
}

func (*MyMiddleware) ProcessResponse(request *goi.Request, response *goi.Response) {
}
```

注册：

```go
package middleware

import "goi_example/backend/goi_example"

func init() {
	goi_example.ApiRouter.Use(&MyMiddleware{})
}
```

执行顺序：

- 请求阶段（ProcessRequest）：正序执行；返回 nil 则继续执行，否则直接写入响应返回
- 异常阶段（ProcessException）：逆序执行；返回 nil 则继续执行，否则直接写入响应返回
- 响应阶段（ProcessResponse）：逆序执行；常用于追加/修改响应头

异常处理：

- 若所有异常中间件均返回 nil，将由框架统一返回 500 响应

## 日志模块

* **日志等级与方法**
    * DEBUG：`goi.Log.Debug()` / `goi.Log.DebugF()`
    * INFO：`goi.Log.Info()` / `goi.Log.InfoF()`
    * WARNING：`goi.Log.Warning()` / `goi.Log.WarningF()`
    * ERROR：`goi.Log.Error()` / `goi.Log.ErrorF()`
    * 自定义等级：`goi.Log.Log()` / `goi.Log.LogF()`
* **Debug 输出**
    * `goi.Settings.Debug` 为 `true` 时，日志同时输出到控制台和日志文件
    * `goi.Settings.Debug` 为 `false` 时，跳过 DEBUG 日志，其他等级仍写入日志文件
* **日志切割**：两种可同时设置满足任意一种即切割
    * 按照日志大小：`Log.SplitSize`
        * 例：`Log.SplitSize = 1024 * 1024`
    * 按照日期格式：`Log.SplitTime = "2006-01-02"`
        * 例：`Log.SplitTime = "2006-01-02"` 设置日期格式
* **日志关闭**
    * 框架提供 `Logger.Close()`，具体关闭时机由使用方管理
    * 日志切割任务通过 `RegisterOnStartup` 注册时，应在 `RunServer()` 返回后关闭日志
    * 不要在 `OnShutdown` 中关闭默认日志，框架执行回调后仍需要记录数据库关闭和最终停机日志

## 缓存

* 读取时会检查过期项，也可配置过期清理策略
    * `goi.Periodic`：定期抽查过期项，默认策略
    * `goi.Scheduled`：定时调度策略；不承诺到期立即清除，实际时机需按当前版本验证
* 支持内存淘汰策略
    * `goi.NoEviction`（不淘汰已有键，容量不足时抛出错误）
    * `goi.AllKeysRandom`（随机删除-所有键）
    * `goi.AllKeysLRU`（删除最近最少使用-所有键）
    * `goi.AllKeysLFU`（删除最近最不频繁使用-所有键）
    * `goi.VolatileRandom`（随机删除-设置过期时间的键）
    * `goi.VolatileLRU`（删除最近最少使用-设置过期时间的键）
    * `goi.VolatileLFU`（删除最近最不频繁使用-设置过期时间的键）
    * `goi.VolatileTTL`（删除即将过期的键-设置过期时间的键）

## 内置 auth 密码生成与验证

```go
package auth_test

import (
	"fmt"

	"github.com/NeverStopDreamingWang/goi/v3/auth"
)

func ExampleMakePassword() {
	// 原始密码
	password := "goi123456"

	// 生成密码的哈希值
	hashedPassword, err := auth.MakePassword(password)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("原始密码: %s\n", password)

	// 验证密码
	isValid := auth.CheckPassword(password, hashedPassword)
	fmt.Printf("密码验证: %v\n", isValid)

	// Output:
	//
	// 原始密码: goi123456
	// 密码验证: true
}
```

## 内置 Converter 路由转换器

注册路由转换器

* 包含参数验证
* 包含参数解析

功能：

实现 `goi.Converter` 结构体

```go
package goi_test

import (
	"errors"

	"github.com/NeverStopDreamingWang/goi/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// 手机号
var phoneConverter = goi.Converter{
	Regex: `(1[3456789]\d{9})`,
	ToGo:  func(value string) (any, error) { return value, nil },
}

// 邮箱
var emailConverter = goi.Converter{
	Regex: `([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`,
	ToGo:  func(value string) (any, error) { return value, nil },
}

// URL：简化示例，不用于完整的 URL 合法性校验
var urlConverter = goi.Converter{
	Regex: `((?:https?://|www\.)[a-zA-Z0-9][a-zA-Z0-9.-]*(?::[0-9]+)?(?:/[^\s]*)?)`,
	ToGo:  func(value string) (any, error) { return value, nil },
}

// 日期 (YYYY-MM-DD)
var dateConverter = goi.Converter{
	Regex: `(\d{4}-(?:0[1-9]|1[0-2])-(?:0[1-9]|[12]\d|3[01]))`,
	ToGo:  func(value string) (any, error) { return value, nil },
}

// 时间 (HH:MM:SS)
var timeConverter = goi.Converter{
	Regex: `((?:[01]\d|2[0-3]):[0-5]\d:[0-5]\d)`,
	ToGo:  func(value string) (any, error) { return value, nil },
}

// IP地址 (IPv4)
var ipv4Converter = goi.Converter{
	Regex: `((?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))`,
	ToGo:  func(value string) (any, error) { return value, nil },
}

// 用户名 (字母开头，允许字母数字下划位)
var usernameConverter = goi.Converter{
	Regex: `([a-zA-Z][a-zA-Z0-9_]{3,15})`,
	ToGo:  func(value string) (any, error) { return value, nil },
}

// ObjectId
var objectIdConverter = goi.Converter{
	Regex: `([a-fA-F0-9]{24})`,
	ToGo: func(value string) (interface{}, error) {
		objectId, err := bson.ObjectIDFromHex(value)
		if err != nil {
			return nil, errors.New("ID 错误")
		}
		return objectId, nil
	},
}

func ExampleRegisterConverter() {
	// 注册路由转换器

	// 手机号
	goi.RegisterConverter("my_phone", phoneConverter)

	// 邮箱
	goi.RegisterConverter("my_email", emailConverter)

	// URL
	goi.RegisterConverter("my_url", urlConverter)

	// 日期 (YYYY-MM-DD)
	goi.RegisterConverter("my_date", dateConverter)

	// 时间 (HH:MM:SS)
	goi.RegisterConverter("my_time", timeConverter)

	// IP地址 (IPv4)
	goi.RegisterConverter("my_ipv4", ipv4Converter)

	// 用户名 (字母开头，允许字母数字下划位)
	goi.RegisterConverter("my_username", usernameConverter)

	// MongoDB objectId
	goi.RegisterConverter("object_id", objectIdConverter)
}
```

使用

```go
package goi_test

import (
	"net/http"

	"goi_example/backend/goi_example"

	"github.com/NeverStopDreamingWang/goi/v3"
)

func init() {
	// 创建一个子路由
	testRouter := goi_example.Server.Router.Include("test/", "测试路由")
	{
		// 注册一个路径 
		// 类型 my_phone
		// 名称 phone
		testRouter.Path("test_phone/<my_phone:phone>", "", goi.ViewSet{GET: TestPhone}) // 测试路由转换器
	}
}

// 测试手机号路由转换器
func TestPhone(request *goi.Request) any {
	var phone string
	var validationErr goi.ValidationError
	validationErr = request.PathParams.Get("phone", &phone)
	if validationErr != nil {
		return validationErr.Response()
	}
	resp := map[string]any{
		"status": http.StatusOK,
		"msg":    phone,
		"data":   "OK",
	}
	return resp
}

```

## 内置 Validator 参数验证器

Goi v3.0.3 的 `goi.Validator` 只需实现 `Parse(value any) (any, goi.ValidationError)`，一次完成验证和转换：成功返回可写入目标字段的 Go 值，失败返回 `nil` 和验证错误。旧自定义验证器需将 `Validate`、`ToGo` 的逻辑合并到 `Parse`；上面的路由转换器 `Converter.ToGo` 不受影响。

当前内置规则：

- `bool` 字符串遵循 `strconv.ParseBool` 的固定格式，例如 `TRUE`、`true`、`1`；不额外转换大小写或去除空白。
- `int` 字符串通过 `strconv.Atoi` 转换，支持正负号，格式错误或超出当前平台 int 范围时返回验证错误；原生数值保持原类型。
- `string` 使用 `fmt.Sprint` 转换，不等同于 JSON 编码；`slug`、`uuid` 必须整串匹配，UUID 仍只接受小写十六进制。
- `slice` 接受切片及其单层指针，不接受原生数组；`slice`、`map` 的字符串直接进行 JSON 解析，兼容外层空白、多行和 `null`。JSON `null` 返回带类型的 nil 集合，不保证清空已有字段。

`ParseParams` 根据字段标签选择验证器，再通过 `SetValue` 写入字段；`Params.Get` 不调用注册的验证器。原生 nil 仍按 `required`、`allow_null` 标签处理，允许 nil 时跳过赋值。

完整规则及自定义方式见 [参数验证与类型转换指南](https://github.com/NeverStopDreamingWang/goi/blob/v3.0.3/docs/validation.md)（[English](https://github.com/NeverStopDreamingWang/goi/blob/v3.0.3/docs/validation.en.md)）。下面的示例分别返回 string、`bson.ObjectID` 和自定义 `Email`：

```go
package goi_test

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"testing"

	"github.com/NeverStopDreamingWang/goi/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// phoneValidator 手机号验证器示例
type phoneValidator struct{}

func (validator phoneValidator) Parse(value any) (any, goi.ValidationError) {
	var valueStr string
	switch typeValue := value.(type) {
	case int:
		valueStr = strconv.Itoa(typeValue)
	case string:
		valueStr = typeValue
	default:
		return nil, goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数类型错误：%v", value))
	}
	if !regexp.MustCompile(`^(1[3456789]\d{9})$`).MatchString(valueStr) {
		return nil, goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数错误：%v", value))
	}
	return valueStr, nil
}

// MongoDB objectId  类型
type objectIdValidator struct{}

func (validator objectIdValidator) Parse(value any) (any, goi.ValidationError) {
	valueStr, ok := value.(string)
	if !ok {
		return nil, goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数类型错误：%v", value))
	}
	objectId, err := bson.ObjectIDFromHex(valueStr)
	if err != nil {
		return nil, goi.NewValidationError(http.StatusBadRequest, "ID 错误")
	}
	return objectId, nil
}

// emailValidator 邮箱验证器示例
type emailValidator struct{}

func (validator emailValidator) Parse(value any) (any, goi.ValidationError) {
	switch v := value.(type) {
	case string:
		re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !re.MatchString(v) {
			return nil, goi.NewValidationError(http.StatusBadRequest, "invalid email format")
		}
		return Email{Address: v}, nil
	default:
		return nil, goi.NewValidationError(http.StatusBadRequest, "email must be string")
	}
}

// 支持自定义类型
type Email struct {
	Address string `json:"address"`
}
type testParamsValidParams struct {
	Phone string `name:"phone" type:"phone" required:"true"`
	Email *Email `name:"email" type:"email" required:"true"`
}

// ExampleRegisterValidator 展示如何注册和使用自定义验证器
func ExampleRegisterValidator() {
	// 注册手机号验证器
	goi.RegisterValidator("phone", phoneValidator{})

	// MongoDB objectId
	goi.RegisterValidator("object_id", objectIdValidator{})

	// 注册邮箱验证器
	goi.RegisterValidator("email", emailValidator{})

	var validationErr goi.ValidationError

	var params testParamsValidParams
	bodyParams := goi.Params{
		"phone": "13800000000",
		"email": "test@example.com",
	}
	validationErr = bodyParams.ParseParams(&params)
	if validationErr != nil {
		return
	}
	fmt.Printf("Phone %+v %T\n", params.Phone, params.Phone)
	fmt.Printf("Email %+v %T\n", params.Email, params.Email)
}

func TestRegisterValidator(t *testing.T) {
	ExampleRegisterValidator()
}

```

## ORM 模型

本仓库演示 MySQL 和 SQLite3，下面以 SQLite3 为例。Goi 还提供其它数据库方言，具体接口与扩展方式见 Goi 仓库的 `db/README.md`。

v3.0.3 的原生方言默认解析器支持在 `ExecuteContext`、`QueryContext`、`QueryRowContext` 及其无 Context 方法中使用 `IN ?`、`NOT IN ?`，集合参数传切片、数组或其非 nil 指针；字节切片不展开，空集合生成 `()`。`IN (?)` 不自动展开，须手动提供括号内的占位符与参数。自定义 `PlaceholderParser.Parse` 现需返回 `(string, []any)`，分别作为最终执行的 SQL 和参数；本示例使用内置解析器，无需自定义实现。

v3 的迁移、CRUD、统计、分页和事务方法直接接收 `context.Context`；HTTP 业务应传入 `request.Object.Context()`，下面的独立示例使用 `context.Background()`。`Insert`、`Update` 还需显式传入 `skipNil`：`true` 跳过 nil，`false` 将选中的 nil 写为 SQL NULL；普通零值不会被跳过，用 `Fields` 控制写入范围。`Update`、`Delete` 无 WHERE 时会操作整张表，务必检查条件。

`Select` 必须传入切片指针，`First` 使用结构体指针或已初始化的直接 map；map 的键是 Go 字段名，不是数据库列名。以下示例共享测试表，输出仅作示意，取决于数据和执行顺序；初始化中的 `DROP TABLE` 会删除测试表，请勿在业务数据库中运行。

```go
package sqlite3_test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/NeverStopDreamingWang/goi/v3"
	"github.com/NeverStopDreamingWang/goi/v3/db"
	"github.com/NeverStopDreamingWang/goi/v3/db/sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	ctx := context.Background()
	goi.Settings.Databases["default"] = &goi.Database{
		Engine: "sqlite3",
		Connect: func(engine string) *sql.DB {
			var DB *sql.DB
			var err error
			DataSourceName := filepath.Join(goi.Settings.BaseDir, "sqlite3.db")
			DB, err = sql.Open(engine, DataSourceName)
			if err != nil {
				goi.Log.Error(err)
				panic(err)
			}
			return DB
		},
	}
	// 初始化测试数据库
	// 使用泛型通用连接
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 使用专用连接
	// sqliteDB := sqlite3.Connect("default")

	// 迁移模型前清理旧表(如果存在)
	_, _ = sqliteDB.ExecuteContext(ctx, "DROP TABLE IF EXISTS main.user_tb")

	// 迁移模型
	sqliteDB.Migrate(ctx, UserModel{})

	// 插入初始测试数据
	username := "test_user"
	password := "test123456"
	now := time.Now().Format(time.DateTime)

	user := UserModel{
		Username:   &username,
		Password:   &password,
		CreateTime: &now,
		UpdateTime: &now,
	}

	sqliteDB.SetModel(UserModel{})
	result, err := sqliteDB.Fields("Username", "Password", "CreateTime", "UpdateTime").Insert(ctx, user, true)
	if err != nil {
		fmt.Println("插入错误:", err)
		return
	}
	id, _ := result.LastInsertId()
	user.Id = &id
}

// 用户表模型
type UserModel struct {
	Id         *int64  `field_name:"id" field_type:"INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT" json:"id"`
	Username   *string `field_name:"username" field_type:"TEXT NOT NULL" json:"username"`
	Password   *string `field_name:"password" field_type:"TEXT NOT NULL" json:"-"`
	CreateTime *string `field_name:"create_time" field_type:"TEXT NOT NULL" json:"create_time"`
	UpdateTime *string `field_name:"update_time" field_type:"TEXT" json:"update_time"`
}

// 设置表配置
func (userModel UserModel) ModelSet() *sqlite3.Settings {
	encryptFields := []string{
		"username",
		"password",
	}

	modelSettings := &sqlite3.Settings{
		MigrationsHandler: sqlite3.MigrationsHandler{
			BeforeHandler: nil,      // 迁移之前处理函数
			AfterHandler:  InitUser, // 迁移之后处理函数
		},

		Schema:    "main",    // SQLite 命名空间；附加数据库需预先 ATTACH
		TableName: "user_tb", // 设置表名

		// 自定义配置
		Settings: goi.Params{
			"encrypt_fields": encryptFields,
		},
	}

	return modelSettings
}

// 初始化用户
func InitUser() error {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	var (
		id              int64 = 1
		username              = "超级管理员"
		password              = "admin"
		create_Datetime       = goi.GetTime().Format(time.DateTime)
		err             error
	)

	user := UserModel{
		Id:         &id,
		Username:   &username,
		Password:   &password,
		CreateTime: &create_Datetime,
		UpdateTime: nil,
	}
	sqliteDB.SetModel(UserModel{})
	_, err = sqliteDB.Insert(ctx, &user, false)
	if err != nil {
		return err
	}

	return nil
}

func ExampleEngine_Migrate() {
	ctx := context.Background()
	// 连接数据库
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 迁移模型
	// Schema 在 UserModel.ModelSet 返回的 Settings.Schema 中指定，不再作为 Migrate 参数。
	sqliteDB.Migrate(ctx, UserModel{})

	// Output:
}

func ExampleEngine_Insert() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 准备测试数据
	username := "test_user"
	password := "test123456"
	now := time.Now().Format(time.DateTime)

	user := UserModel{
		Username:   &username,
		Password:   &password,
		CreateTime: &now,
	}

	// 设置模型并插入数据
	sqliteDB.SetModel(UserModel{})
	result, err := sqliteDB.Fields("Username", "Password", "CreateTime").Insert(ctx, user, true)
	if err != nil {
		fmt.Println("插入错误:", err)
		return
	}

	id, _ := result.LastInsertId()
	fmt.Printf("插入成功, ID: %d\n", id)

	// Output:
	// 插入成功, ID: 1
}

func ExampleEngine_Where() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	var users []*UserModel
	sqliteDB.SetModel(UserModel{})

	// Where 返回独立条件副本，需接收返回值或继续链式调用
	// 简单条件
	sqliteDB = sqliteDB.Where("username = ?", "test_user")

	// AND条件
	sqliteDB = sqliteDB.Where("create_time IS NOT NULL")

	// 支持 IN 参数值为 Slice 或 Array
	testUsers := []string{"test_user", "test_user_1"}
	sqliteDB = sqliteDB.Where("username IN ?", testUsers)

	// LIKE条件
	sqliteDB = sqliteDB.Where("username LIKE ?", "test_%")

	err := sqliteDB.Select(ctx, &users)
	if err != nil {
		fmt.Println("查询错误:", err)
		return
	}

	fmt.Printf("查询到 %d 条记录\n", len(users))
	for _, user := range users {
		fmt.Printf("用户名: %s\n", *user.Username)
	}

	// Output:
	// 查询到 1 条记录
	// 用户名: test_user
}

func ExampleEngine_Select() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 查询多条记录 Select 支持 map 以及 结构体
	// var users []map[string]any
	var users []*UserModel

	sqliteDB.SetModel(UserModel{})
	err := sqliteDB.Select(ctx, &users)
	if err != nil {
		fmt.Println("查询错误:", err)
		return
	}

	for _, user := range users {
		fmt.Printf("用户名: %s\n", *user.Username)
	}

	// Output:
	// 用户名: test_user
}

func ExampleEngine_Page() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 分页查询
	var users []*UserModel

	sqliteDB.SetModel(UserModel{})
	total, totalPages, err := sqliteDB.Page(ctx, 2, 5) // 第2页,每页5条
	if err != nil {
		fmt.Println("分页错误:", err)
		return
	}

	err = sqliteDB.Select(ctx, &users)
	if err != nil {
		fmt.Println("查询错误:", err)
		return
	}

	fmt.Printf("总记录数: %d\n", total)
	fmt.Printf("总页数: %d\n", totalPages)
	fmt.Printf("当前页记录数: %d\n", len(users))

	// Output:
	// 总记录数: 15
	// 总页数: 3
	// 当前页记录数: 5
}

func ExampleEngine_First() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 查询单条记录 First 支持 map 以及 结构体
	// 使用 map 时先 make，再把 map 本身传给 First。
	var user UserModel

	sqliteDB.SetModel(UserModel{})
	sqliteDB = sqliteDB.Where("username = ?", "test_user")
	err := sqliteDB.First(ctx, &user)
	if err != nil {
		fmt.Println("查询错误:", err)
		return
	}

	fmt.Printf("用户名: %s\n", *user.Username)

	// Output:
	// 用户名: test_user
}

func ExampleEngine_Fields() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 指定查询字段
	user := make(map[string]*string)

	sqliteDB.SetModel(UserModel{})

	// Fields 返回当前实例的副本指针
	sqliteDB = sqliteDB.Fields("Username", "CreateTime")
	// First 写入单条 map；Select 则需要 *[]map[string]any 等切片指针。
	err := sqliteDB.Where("username = ?", "test_user").First(ctx, user)
	if err != nil {
		fmt.Println("查询错误:", err)
		return
	}

	// map 键使用 Go 字段名；值类型由模型字段决定，这里是 *string。
	fmt.Println("用户名:", *user["Username"])

	// Output:
	// 用户名: test_user
}

func ExampleEngine_Update() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 准备更新数据
	newPassword := "new_password"
	updateUser := UserModel{
		Password: &newPassword,
	}

	// 更新记录
	sqliteDB.SetModel(UserModel{})
	sqliteDB = sqliteDB.Where("username = ?", "test_user")
	result, err := sqliteDB.Fields("Password").Update(ctx, updateUser, false)
	if err != nil {
		fmt.Println("更新错误:", err)
		return
	}

	rows, _ := result.RowsAffected()
	fmt.Printf("更新成功, 影响行数: %d\n", rows)

	// Output:
	// 更新成功, 影响行数: 1
}

func ExampleEngine_Delete() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	sqliteDB.SetModel(UserModel{})
	sqliteDB = sqliteDB.Where("username = ?", "test_user")

	// 删除记录
	result, err := sqliteDB.Delete(ctx)
	if err != nil {
		fmt.Println("删除错误:", err)
		return
	}

	rows, _ := result.RowsAffected()
	fmt.Printf("删除成功, 影响行数: %d\n", rows)

	// Output:
	// 删除成功, 影响行数: 1
}

func ExampleEngine_GroupBy() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 使用分组查询
	type Result struct {
		CreateTime *string `json:"create_time"`
	}

	var results []*Result
	sqliteDB.SetModel(UserModel{})
	// GroupBy 接收列名，Fields 接收 Go 字段名，均返回副本指针。
	sqliteDB = sqliteDB.Fields("CreateTime").GroupBy("create_time")
	err := sqliteDB.Select(ctx, &results)
	if err != nil {
		fmt.Println("查询错误:", err)
		return
	}

	for _, result := range results {
		fmt.Printf("创建时间: %s\n", *result.CreateTime)
	}

	// Output:
	// 创建时间: 2024-03-21 10:00:00
}

func ExampleEngine_OrderBy() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 使用排序查询
	var users []UserModel
	sqliteDB.SetModel(UserModel{})

	// OrderBy 返回当前实例的副本指针
	sqliteDB = sqliteDB.OrderBy("-create_time")
	err := sqliteDB.Select(ctx, &users)
	if err != nil {
		fmt.Println("查询错误:", err)
		return
	}

	for _, user := range users {
		fmt.Println("用户名:", *user.Username)
	}

	// Output:
	// 用户名: test_user
}

func ExampleEngine_Count() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 统计记录数
	sqliteDB.SetModel(UserModel{})
	count, err := sqliteDB.Count(ctx)
	if err != nil {
		fmt.Println("统计错误:", err)
		return
	}

	fmt.Printf("总记录数: %d\n", count)

	// Output:
	// 总记录数: 1
}

func ExampleEngine_Exists() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	sqliteDB.SetModel(UserModel{})
	// 检查记录是否存在
	exists, err := sqliteDB.Where("username = ?", "test_user").Exists(ctx)
	if err != nil {
		fmt.Println("查询错误:", err)
		return
	}

	fmt.Printf("记录是否存在: %v\n", exists)

	// Output:
	// 记录是否存在: true
}

func ExampleEngine_WithTransaction() {
	ctx := context.Background()
	sqliteDB := db.Connect[*sqlite3.Engine]("default")

	// 事务
	err := sqliteDB.WithTransaction(ctx, func(engine *sqlite3.Engine, args ...any) error {
		// 准备测试数据
		username1 := "transaction_user1"
		username2 := "transaction_user2"
		password := "test123456"
		now := time.Now().Format(time.DateTime)

		// 插入第一个用户
		user1 := UserModel{
			Username:   &username1,
			Password:   &password,
			CreateTime: &now,
		}
		engine.SetModel(UserModel{})
		_, err := engine.Fields("Username", "Password", "CreateTime").Insert(ctx, user1, true)
		if err != nil {
			return err
		}

		// 插入第二个用户
		user2 := UserModel{
			Username:   &username2,
			Password:   &password,
			CreateTime: &now,
		}
		_, err = engine.Fields("Username", "Password", "CreateTime").Insert(ctx, user2, true)
		if err != nil {
			return err
		}

		return nil
	})
	// 错误向外返回时由最外层回滚；嵌套 WithTransaction 复用外层事务，不创建保存点。
	if err != nil {
		fmt.Println("事务执行错误:", err)
		return
	}

	fmt.Println("事务执行完成")

	// Output:
	// 事务执行完成
}

```

`GroupBy` 默认把输入当作列标识符引用，不能直接把 `DATE(create_time)` 或 `COUNT(*)` 当作字段传入。复杂表达式和聚合查询请使用 `QueryContext` 编写明确的 SQL，再通过 `Rows.Scan` 读取结果，并关闭结果集、检查 `Rows.Err()`。`Migrate` 只创建不存在的表，不会修改已有表结构；SQLite 的 `Schema` 仅支持 `main`、`temp` 或已通过 ATTACH 建立的数据库名称。
