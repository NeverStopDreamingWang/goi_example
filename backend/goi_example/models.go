package goi_example

import (
	"path"

	"goi_example/backend/utils"
	"goi_example/backend/utils/mongodb"
	"goi_example/backend/utils/mysqldb"
	"goi_example/backend/utils/redisdb"
	"goi_example/backend/utils/sqlite3db"
)

// 配置
type ConfigModel struct {
	Debug            bool                   `yaml:"debug"`
	Port             uint16                 `yaml:"port"`
	CorsAllowOrigins []string               `yaml:"cors_allow_origins"`
	SQLite3Config    *sqlite3db.ConfigModel `yaml:"sqlite3"`
	MySQLConfig      *mysqldb.ConfigModel   `yaml:"mysql"`
	RedisConfig      *redisdb.ConfigModel   `yaml:"redis"`
	MongoDBConfig    *mongodb.ConfigModel   `yaml:"mongodb"`
}

// 加载配置
func (self *ConfigModel) Load() error {
	configPath := path.Join(Server.Settings.BASE_DIR, "config.yaml")
	return utils.LoadYaml(configPath, self)
}

// 保存配置
func (self *ConfigModel) Save() error {
	configPath := path.Join(Server.Settings.BASE_DIR, "config.yaml")
	return utils.SaveYaml(configPath, self)
}
