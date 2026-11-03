package goi_example

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/NeverStopDreamingWang/goi/v3"
	"github.com/NeverStopDreamingWang/goi/v3/response"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func init() {
	// 注册验证器
	// 手机号
	goi.RegisterValidator("phone", phoneValidator{})
	// 注册验证器
	goi.RegisterValidator("object_id", objectIdValidator{})
}

// 自定义参数验证错误
type validationError struct {
	Status  int
	Message string
}

// 创建参数验证错误方法
func (validationErr validationError) NewValidationError(status int, message string, args ...any) goi.ValidationError {
	return &validationError{
		Status:  status,
		Message: message,
	}
}

// 实现 error 接口，返回错误消息
func (validationErr validationError) Error() string {
	return validationErr.Message
}

// 参数验证错误响应格式
func (validationErr validationError) Response() goi.Response {
	return goi.Response{
		Status: http.StatusOK,
		Data: response.Data{
			Code:    validationErr.Status,
			Message: validationErr.Message,
			Data:    nil,
		},
	}
}

type phoneValidator struct{}

var phoneRe = regexp.MustCompile(`^(1[3456789]\d{9})$`)

// Parse 验证手机号并返回字符串，与参数结构体的手机号字段类型保持一致。
func (validator phoneValidator) Parse(value any) (any, goi.ValidationError) {
	var phone string
	switch typeValue := value.(type) {
	case int:
		phone = strconv.Itoa(typeValue)
	case string:
		phone = typeValue
	default:
		return nil, goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数类型错误：%v", value))
	}
	if !phoneRe.MatchString(phone) {
		return nil, goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数错误：%v", value))
	}
	return phone, nil
}

// MongoDB objectId  类型
type objectIdValidator struct{}

// Parse 验证十六进制 ObjectID 字符串并转换为 bson.ObjectID。
func (validator objectIdValidator) Parse(value any) (any, goi.ValidationError) {
	text, ok := value.(string)
	if !ok {
		return nil, goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数类型错误：%v", value))
	}
	objectId, err := bson.ObjectIDFromHex(text)
	if err != nil {
		return nil, goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数错误：%v", value))
	}
	return objectId, nil
}
