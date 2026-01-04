package goi_example

import (
	"path"

	"goi_example/backend/utils"
	"goi_example/backend/utils/mongo_db"
	"goi_example/backend/utils/mysql_db"
	"goi_example/backend/utils/redis_db"
	"goi_example/backend/utils/sqlite3_db"
)

// 配置
type ConfigModel struct {
	Debug            bool                    `yaml:"debug"`
	Port             uint16                  `yaml:"port"`
	CorsAllowOrigins []string                `yaml:"cors_allow_origins"`
	SQLite3Config    *sqlite3_db.ConfigModel `yaml:"sqlite3"`
	MySQLConfig      *mysql_db.ConfigModel   `yaml:"mysql"`
	RedisConfig      *redis_db.ConfigModel   `yaml:"redis"`
	MongoDBConfig    *mongo_db.ConfigModel   `yaml:"mongodb"`
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
