package goi_example

import (
	"errors"

	"github.com/NeverStopDreamingWang/goi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	// 注册路由转换器
	// 手机号
	goi.RegisterConverter("phone", phoneConverter)
	// ObjectId
	goi.RegisterConverter("object_id", objectIdConverter)
}

// 手机号
var phoneConverter = goi.Converter{
	Regex: `(1[3456789]\d{9})`,
	ToGo:  func(value string) (interface{}, error) { return value, nil },
}

// ObjectId
var objectIdConverter = goi.Converter{
	Regex: `([a-fA-F0-9]{24})`,
	ToGo: func(value string) (interface{}, error) {
		objectId, err := primitive.ObjectIDFromHex(value)
		if err != nil {
			return nil, errors.New("ID 错误")
		}
		return objectId, nil
	},
}
