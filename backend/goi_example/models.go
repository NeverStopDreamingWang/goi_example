package goi_example

import (
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
