package mysql_db

import (
	"database/sql"
	"time"

	"github.com/NeverStopDreamingWang/goi"
	_ "github.com/go-sql-driver/mysql"
)

// MySQL 配置
type ConfigModel struct {
	Uri string `yaml:"uri"`
}

var Config *ConfigModel

func Connect(ENGINE string) *sql.DB {
	if Config == nil {
		panic("MySQL 配置未初始化")
	}
	var err error
	sqlite3DB, err := sql.Open(ENGINE, Config.Uri)
	if err != nil {
		goi.Log.Error(err)
		panic(err)
	}
	// 设置连接池参数
	sqlite3DB.SetMaxOpenConns(10)           // 设置最大打开连接数
	sqlite3DB.SetMaxIdleConns(5)            // 设置最大空闲连接数
	sqlite3DB.SetConnMaxLifetime(time.Hour) // 设置连接的最大存活时间
	return sqlite3DB
}
