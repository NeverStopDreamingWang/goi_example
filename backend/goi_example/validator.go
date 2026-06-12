package goi_example

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/NeverStopDreamingWang/goi/v2"
	"github.com/NeverStopDreamingWang/goi/v2/response"
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

func (validator phoneValidator) Validate(value any) goi.ValidationError {
	switch typeValue := value.(type) {
	case int:
		valueStr := strconv.Itoa(typeValue)
		var reStr = `^(1[3456789]\d{9})$`
		re := regexp.MustCompile(reStr)
		if re.MatchString(valueStr) == false {
			return goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数错误：%v", value))
		}
	case string:
		var reStr = `^(1[3456789]\d{9})$`
		re := regexp.MustCompile(reStr)
		if re.MatchString(typeValue) == false {
			return goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数错误：%v", value))
		}
	default:
		return goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数类型错误：%v", value))
	}
	return nil
}

func (validator phoneValidator) ToGo(value any) (any, goi.ValidationError) {
	switch typeValue := value.(type) {
	case int:
		return typeValue, nil
	case string:
		intValue, err := strconv.Atoi(typeValue)
		if err != nil {
			return nil, goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数类型错误：%v", value))
		}
		return intValue, nil
	default:
		// 尝试转换为字符串
		return fmt.Sprintf("%v", value), nil
	}
}

// MongoDB objectId  类型
type objectIdValidator struct{}

func (validator objectIdValidator) Validate(value any) goi.ValidationError {
	switch typeValue := value.(type) {
	case string:
		var reStr = `^([a-fA-F0-9]{24})$`
		re := regexp.MustCompile(reStr)
		if re.MatchString(typeValue) == false {
			return goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数错误：%v", value))
		}
	default:
		return goi.NewValidationError(http.StatusBadRequest, fmt.Sprintf("参数类型错误：%v", value))
	}
	return nil
}

func (validator objectIdValidator) ToGo(value any) (any, goi.ValidationError) {
	switch typeValue := value.(type) {
	case string:
		objectId, err := bson.ObjectIDFromHex(typeValue)
		if err != nil {
			return nil, goi.NewValidationError(http.StatusBadRequest, "ID 错误")
		}
		return objectId, nil
	default:
		return typeValue, nil
	}
}
