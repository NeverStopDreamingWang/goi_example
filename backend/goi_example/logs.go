package goi_example

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/NeverStopDreamingWang/goi"
)

// 日志输出
func LogPrintln(logger *goi.MetaLogger, level goi.Level, logs ...interface{}) {
	timeStr := fmt.Sprintf("[%v]", goi.GetTime().Format(time.DateTime))
	if level != "" {
		timeStr += fmt.Sprintf(" %v", level)
	}
	logs = append([]interface{}{timeStr}, logs...)
	logger.Logger.Println(logs...)
}

// 获取 *os.File 对象
func getFileFunc(filePath string) (*os.File, error) {
	var file *os.File
	var err error

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) { // 不存在则创建
		file, err = os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		// 文件存在，打开文件
		file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}
	return file, nil
}

// 默认日志
func newDefaultLog() *goi.MetaLogger {
	var err error

	OutPath := filepath.Join(Server.Settings.BASE_DIR, "logs", "server.log")
	OutDir := filepath.Dir(OutPath) // 检查目录
	_, err = os.Stat(OutDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(OutDir, 0755)
		if err != nil {
			panic(fmt.Sprintf("创建日志目录错误: %v", err))
		}
	}

	defaultLog := &goi.MetaLogger{
		Name: "默认日志",
		Path: OutPath,
		Level: []goi.Level{ // 所有等级的日志
			goi.INFO,
			goi.WARNING,
			goi.ERROR,
		},
		Logger:          nil,
		File:            nil,
		LoggerPrint:     LogPrintln, // 日志输出格式
		CreateTime:      goi.GetTime(),
		SPLIT_SIZE:      1024 * 1024 * 2, // 切割大小
		SPLIT_TIME:      "2006-01-02",    // 切割日期，每天
		GetFileFunc:     getFileFunc,     // 创建文件对象方法
		SplitLoggerFunc: nil,             // 自定义日志切割：符合切割条件时，传入日志对象，返回新的文件对象
	}
	defaultLog.File, err = getFileFunc(defaultLog.Path)
	if err != nil {
		panic(fmt.Sprintf("初始化[%v]日志错误: %v", defaultLog.Name, err))
	}
	defaultLog.Logger = log.New(defaultLog.File, "", 0)
	return defaultLog
}
