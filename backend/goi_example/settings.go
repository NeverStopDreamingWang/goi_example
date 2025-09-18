package goi_example

import (
	"os"
	"path"
	"path/filepath"

	"goi_example/backend/utils"
	"goi_example/backend/utils/mongo_db"
	"goi_example/backend/utils/mysql_db"
	"goi_example/backend/utils/redis_db"
	"goi_example/backend/utils/sqlite3_db"

	"github.com/NeverStopDreamingWang/goi"
	"github.com/NeverStopDreamingWang/goi/middleware"
)

// Http 服务
var Server *goi.Engine
var Config *ConfigModel

var STATIC_DIR string
var STATIC_URL = "/static/"

func init() {
	var err error

	// 创建 http 服务
	Server = goi.NewHttpServer()

	// version := goi.Version() // 获取版本信息
	// fmt.Println("goi 版本", version)

	// 注册中间件
	Server.MiddleWare = []goi.MiddleWare{
		&middleware.SecurityMiddleWare{},
		&middleware.CommonMiddleWare{},
		&middleware.XFrameMiddleWare{},
	}

	// 项目路径
	Server.Settings.BASE_DIR, _ = os.Getwd()
	// 网络协议
	Server.Settings.NET_WORK = "tcp" // 默认 "tcp" 常用网络协议 "tcp"、"tcp4"、"tcp6"、"udp"、"udp4"、"udp6
	// 监听地址
	Server.Settings.BIND_ADDRESS = "0.0.0.0" // 默认 127.0.0.1
	// 端口
	Server.Settings.PORT = 8080
	// 域名
	Server.Settings.BIND_DOMAIN = ""

	// 密钥
	Server.Settings.SECRET_KEY = "goi-insecure-_1pnr2e-&esfi965^#@dg0w4a7jhkqn)aype2m0il0z)vsp8b#"

	// RSA 私钥
	Server.Settings.PRIVATE_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAvA7HhHLEuZ/zimilDOr8sDjRMlEfH1XhQyoCSNoSQOfVAK4Z
O42c7ys1TED6EmCAK5CY5j0KWtkFkZKlB43kcmM1Z3uzTeQ/cEj/q6MrdJBcwtmF
/WA5hxrpUfjSkFvIEokTlLhK1kb6b3D5BV+JhTQs45is/pSIOFEDkW7u05xzDyCf
8Y+4WUHgb6BmG3pj6dVXROE7HTbijo1ZIB9N2NbcJWnadLf2wNJaZjfe/oCSiCEx
bOpvLhDc1JKeH9AhL6YL8swVOpMhr2yovGWeYG056vLbp8cbRXgGoZUBvQbXv34G
PmKkon1TtFcoxUbvY6VUtEgWQY6dtqTkT455BwIDAQABAoIBAB3k8454eBFR/fL4
o5QkHrscrRSklJ/0lPSKTwaps9EhiVisVFKFNndGlRhkE9yr/nPubn/bIDRE33++
ogFLaw9L+gdTQLOXHwaIdCwiqhvDfxtyXLxkeYCipIhlV4OfM3TO2ZAAo9TgP5tL
iCp0f3XvT1t2v7lQfz8Ekd6ildCJfmJAcSY5uJN6zWVywhVs7vZZCGW3vqae2K8O
95NU/BL2TxBqtINpseDsDZ06Yh8MYWXJ32Y6x/mhBkJp2McSiEcot9PKtXsvF3EP
yqSwEBzujrpoajg8nUQ+cVcZuTWwwSqzOmDrBp8Pe+5BD6Rl6+WZSHagTRsqa+Jl
3P7QT9ECgYEA0TFVlh13qlBAbNqsjxJRz5DHxCy8x5PbkTIYKU8cqq962IFjfBes
5ce13MNQv3WWFmd1OEV6wANzAIpT949TjSJjtDLmEoAUZNPCG3A7PkL7SDbNIPfs
ULMdRUAnZE2KZCpq+DyghTkNUP+NqdKxpndFNGLbutbxwaTy8vciNCUCgYEA5iLS
kokTg4NYOKZlOwrzsG1/bUhmZVMQJFiVdb5SOAZ04Z0zBOpU/7Uxj6QuQMVzff2y
rVHYNk6dqqqqfwdWgaE2YEnWWd3RFy8ibPKnYeix/UdC9J0cEaoZWki7HDJZiXcv
I7nGQyOAJ7gnXeLE2cR9RKOTe9ELQ6XsvGYaOrsCgYEAg5evKRE8V4zIGjGs3ws9
H38JyyQBVOJz+nAytrmnZM+iTVOHS2ZxQtJQWqEayHWlhk5qdI1wXB1PWIWrsE0e
1+dMJOznwbeEHLEAp9X/znjALXsbqqOKqnEh9pAWt4f3iG8Ofz1UFLoA4HUBnlSF
oBvjEsMlSfEwfwnOMny3rWkCgYA12JP4YUZFkSfFKXmqFOfrsdMM2NHMh2DRgECI
Kh3GqgwS9dsIHWQB6H1OJJYF5a0eH4v87ZdvLXnKguAdlLPy5Kt6YAxdPn87s3WU
lDoBuJZcsp3B6ji1EV2ZOEc/U7CLb22CKGdxMg88O+RKHVL9uPGua6+IWuMN0vbP
JfyhHQKBgGp/yC8O0a4yFfHSg+Azka04ZTLyf6sH11HtGFuSrdjr7bwSD1YW1xhq
dYuZQ2zUUKQvBfpfXrKEslSPL0yTbDpeeWu+qK++kzSHqGQesKVY71grB6+F/NTW
UDoMYHJTxKqxjolrxYqDbZhPmGQv88AGPp6hmhQORkbPSeLEaKwO
-----END RSA PRIVATE KEY-----
`

	// RSA 公钥
	Server.Settings.PUBLIC_KEY = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvA7HhHLEuZ/zimilDOr8
sDjRMlEfH1XhQyoCSNoSQOfVAK4ZO42c7ys1TED6EmCAK5CY5j0KWtkFkZKlB43k
cmM1Z3uzTeQ/cEj/q6MrdJBcwtmF/WA5hxrpUfjSkFvIEokTlLhK1kb6b3D5BV+J
hTQs45is/pSIOFEDkW7u05xzDyCf8Y+4WUHgb6BmG3pj6dVXROE7HTbijo1ZIB9N
2NbcJWnadLf2wNJaZjfe/oCSiCExbOpvLhDc1JKeH9AhL6YL8swVOpMhr2yovGWe
YG056vLbp8cbRXgGoZUBvQbXv34GPmKkon1TtFcoxUbvY6VUtEgWQY6dtqTkT455
BwIDAQAB
-----END RSA PUBLIC KEY-----
`

	// 设置 SSL
	Server.Settings.SSL = goi.MetaSSL{
		STATUS:    false,  // SSL 开关
		TYPE:      "自签证书", // 证书类型
		CERT_PATH: filepath.Join(Server.Settings.BASE_DIR, "ssl", "goi_example.crt"),
		KEY_PATH:  filepath.Join(Server.Settings.BASE_DIR, "ssl", "goi_example.key"),
	}

	// 加载 Config 配置
	configPath := path.Join(Server.Settings.BASE_DIR, "config.yaml")

	Config = &ConfigModel{}
	err = utils.LoadYaml(configPath, Config)
	if err != nil {
		panic(err)
	}

	if Config.SQLite3Config != nil {
		sqlite3_db.Config = Config.SQLite3Config
		// 数据库配置
		Server.Settings.DATABASES["default"] = &goi.DataBase{
			ENGINE:  "sqlite3",
			Connect: sqlite3_db.Connect,
		}
	}

	if Config.MySQLConfig != nil {
		mysql_db.Config = Config.MySQLConfig
		// 数据库配置
		Server.Settings.DATABASES["mysql"] = &goi.DataBase{
			ENGINE:  "mysql",
			Connect: mysql_db.Connect,
		}
	}

	Server.Settings.USE_TZ = false
	// 设置时区
	err = Server.Settings.SetTimeZone("Asia/Shanghai") // 默认为空字符串 ''，本地时间
	if err != nil {
		panic(err)
	}
	//  goi.GetLocation() 获取时区 Location
	//  goi.GetTime() 获取当前时区的时间

	// 设置框架语言
	Server.Settings.SetLanguage(goi.ZH_CN) // 默认 ZH_CN

	// 设置最大缓存大小
	Server.Cache.EVICT_POLICY = goi.ALLKEYS_LRU   // 缓存淘汰策略
	Server.Cache.EXPIRATION_POLICY = goi.PERIODIC // 过期策略
	Server.Cache.MAX_SIZE = 1024 * 1024 * 20      // 单位为字节，0 为不限制使用

	// 日志 DEBUG 设置
	// 日志 DEBUG 设置
	Server.Log.DEBUG = false
	if Config.Debug != nil {
		Server.Log.DEBUG = *Config.Debug
	}
	// 注册日志
	defaultLog := newDefaultLog() // 默认日志
	err = Server.Log.RegisterLogger(defaultLog)
	if err != nil {
		panic(err)
	}

	STATIC_DIR = path.Join(Server.Settings.BASE_DIR, "static/")
	err = os.MkdirAll(STATIC_DIR, os.ModePerm)
	if err != nil {
		goi.Log.Error(err)
		panic(err)
	}
	// 设置验证器错误，不指定则使用默认
	Server.Validator.SetValidationError(validationError{})

	// 设置自定义配置
	// Server.Settings.Set(key string, value interface{})
	// Server.Settings.Get(key string, dest interface{})

	if Config.RedisConfig != nil {
		redis_db.Config = Config.RedisConfig
		redis_db.Connect()
	}
	if Config.MongoDBConfig != nil {
		mongo_db.Config = Config.MongoDBConfig
		mongo_db.Connect()
	}

	// 注册关闭回调处理程序
	Server.RegisterShutdownHandler("关闭操作", Shutdown)
}

func Shutdown(engine *goi.Engine) error {
	var err error

	if redis_db.Config != nil {
		err = redis_db.Close()
		if err != nil {
			return err
		}
	}

	if mongo_db.Config != nil {
		err = mongo_db.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
