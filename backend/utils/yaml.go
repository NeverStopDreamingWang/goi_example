package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadYaml(filePath string, value any) error {
	// 打开 YAML 文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 YAML 解码器读取文件内容
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(value)
	if err != nil {
		return err
	}
	return nil
}

func SaveYaml(filePath string, value any) error {
	// 打开 YAML 文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 YAML 编码器将数据写入文件
	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(value)
	if err != nil {
		return err
	}
	return nil
}
