package goi_example

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/NeverStopDreamingWang/goi"
)

// 默认日志
func newDefaultLog() *goi.Logger {
	// 使用默认
	// return goi.NewLogger(filepath.Join(Server.Settings.BASE_DIR, "logs", "server.log"))

	var err error

	path := filepath.Join(Server.Settings.BASE_DIR, "logs", "server.log")
	OutDir := filepath.Dir(path) // 检查目录
	_, err = os.Stat(OutDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(OutDir, 0755)
		if err != nil {
			panic(fmt.Sprintf("创建日志目录错误: %v", err))
		}
	}

	logger := &goi.Logger{
		Name:       "默认日志",
		Path:       path,
		File:       nil,
		Logger:     nil,
		SplitSize:  1024 * 1024 * 10, // 切割大小
		SplitTime:  "2006-01-02",     // 切割日期，每天
		CreateTime: goi.GetTime(),
		PrintFunc:  nil, // 自定义日志输出格式，为 nil 则使用默认方法
		// PrintFunc:       printlnFunc,
		GetFileFunc: nil, // 创建文件对象方法 为 nil 则使用默认方法
		// GetFileFunc:     getFileFunc,
		SplitLoggerFunc: nil, // 自定义日志切割：符合切割条件时，传入日志对象，返回新的文件对象
	}

	logger.File, err = logger.GetFile()
	if err != nil {
		panic(fmt.Sprintf("初始化[%v]日志错误: %v", logger.Name, err))
	}
	logger.Logger = log.New(logger.File, "", 0)
	return logger
}

// 日志输出
func printlnFunc(logger *goi.Logger, level goi.Level, logs ...any) {
	timeStr := fmt.Sprintf("[%v]", goi.GetTime().Format(time.DateTime))
	if level != "" {
		timeStr += fmt.Sprintf(" %v", level)
	}
	logs = append([]any{timeStr}, logs...)
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
