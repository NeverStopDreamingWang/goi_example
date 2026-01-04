package sqlite3_db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/NeverStopDreamingWang/goi"
	_ "github.com/mattn/go-sqlite3"
)

// SQLite3 配置
type ConfigModel struct {
	Uri string `yaml:"uri"`
}

var Config *ConfigModel

var SQLite3DB *sql.DB

func Connect(ENGINE string) *sql.DB {
	if Config == nil {
		panic("SQLite3 配置未初始化")
	}
	var err error
	DB_DIR := filepath.Join(goi.Settings.BASE_DIR, "db/")
	_, err = os.Stat(DB_DIR)
	if os.IsNotExist(err) {
		err = os.MkdirAll(DB_DIR, 0755)
		if err != nil {
			panic(fmt.Sprintf("创建数据库目录错误: ", err))
		}
	}
	DataSourceName := filepath.Join(DB_DIR, Config.Uri)
	SQLite3DB, err = sql.Open(ENGINE, DataSourceName)
	if err != nil {
		goi.Log.Error(err)
		panic(err)
	}
	return SQLite3DB
}
