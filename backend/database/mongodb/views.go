package mongodb

import (
	"errors"
	"net/http"

	"goi_example/backend/utils/mongodb"

	"github.com/NeverStopDreamingWang/goi/v2"
	"github.com/NeverStopDreamingWang/goi/v2/response"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// 参数验证
type listValidParams struct {
	Page     int64   `name:"page" type:"int" required:"true"`
	PageSize int64   `name:"page_size" type:"int" required:"true"`
	Search   *string `name:"search" type:"string"`
}

func listView(request *goi.Request) any {
	var params listValidParams
	var queryParams goi.Params
	var validationErr goi.ValidationError

	queryParams = request.QueryParams() // Query 传参
	validationErr = queryParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	collection := DocumentModel{}.Collection()

	// 计算skip值
	skip := (params.Page - 1) * params.PageSize

	// 设置分页查询选项
	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(params.PageSize)

	// 构建查询条件
	filter := bson.M{}
	if params.Search != nil {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": *params.Search}},
			{"content": bson.M{"$regex": *params.Search}},
		}
	}

	ctx, cancel := mongodb.WithTimeout(10)
	defer cancel()
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询标准失败",
			Data:    err.Error(),
		}
	}
	defer cursor.Close(ctx)

	var documentList []*DocumentModel
	if err = cursor.All(ctx, &documentList); err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询标准失败",
			Data:    err.Error(),
		}
	}

	// 获取总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		total = 0
	}

	return response.Data{
		Code:    http.StatusOK,
		Message: "",
		Data: map[string]any{
			"list":  documentList,
			"total": total,
			"page":  params.Page,
			"size":  params.PageSize,
		},
	}
}

// 参数验证
type createValidParams struct {
	Name    *string `name:"name" type:"string" required:"true"`
	Content *string `name:"content" type:"string" required:"true"`
}

func createView(request *goi.Request) any {
	var params createValidParams
	var bodyParams goi.Params
	var validationErr goi.ValidationError

	bodyParams = request.BodyParams() // Body 传参
	validationErr = bodyParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	ctx, cancel := mongodb.WithTimeout(5)
	defer cancel()

	document := &DocumentModel{
		Name:    params.Name,
		Content: params.Content,
	}

	// 参数验证
	err := document.Validate(ctx)
	if err != nil {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}

	// 创建
	err = document.Create(ctx)
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return response.Data{
		Code:    http.StatusOK,
		Message: "创建成功",
		Data:    document,
	}
}

func retrieveView(request *goi.Request) any {
	var pk bson.ObjectID // object_id 转换器将字符串转换为 bson.ObjectID 类型
	var validationErr goi.ValidationError
	validationErr = request.PathParams.Get("pk", &pk) // 路由转换器自动转换
	if validationErr != nil {
		return validationErr.Response()
	}

	// 执行查询操作
	instance := &DocumentModel{}

	collection := instance.Collection()

	filter := bson.M{"_id": pk}

	ctx, cancel := mongodb.WithTimeout(5)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(instance)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return response.Data{
				Code:    http.StatusNotFound,
				Message: "数据不存在",
				Data:    nil,
			}
		}
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询失败",
			Data:    err.Error(),
		}
	}

	return response.Data{
		Code:    http.StatusOK,
		Message: "",
		Data:    instance,
	}
}

// 参数验证
type updateValidParams struct {
	Name    *string `name:"name" type:"string"`
	Content *string `name:"content" type:"string"`
}

func updateView(request *goi.Request) any {
	var pk bson.ObjectID
	var params updateValidParams
	var bodyParams goi.Params
	var validationErr goi.ValidationError

	validationErr = request.PathParams.Get("pk", &pk)
	if validationErr != nil {
		return validationErr.Response()
	}

	bodyParams = request.BodyParams() // Body 传参
	validationErr = bodyParams.ParseParams(&params)
	if validationErr != nil {
		return validationErr.Response()
	}

	ctx, cancel := mongodb.WithTimeout(5)
	defer cancel()

	// 执行查询操作
	instance := &DocumentModel{}

	collection := instance.Collection()

	filter := bson.M{"_id": pk}

	err := collection.FindOne(ctx, filter).Decode(instance)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return response.Data{
				Code:    http.StatusNotFound,
				Message: "数据不存在",
				Data:    nil,
			}
		}
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询失败",
			Data:    err.Error(),
		}
	}

	document := &DocumentModel{
		Id:      instance.Id,
		Name:    params.Name,
		Content: params.Content,
	}
	err = document.Validate(ctx)
	if err != nil {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}
	err = instance.Update(ctx, document)
	if err != nil {
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return response.Data{
		Code:    http.StatusOK,
		Message: "修改成功",
		Data:    instance,
	}
}

func deleteView(request *goi.Request) any {
	var pk bson.ObjectID
	var validationErr goi.ValidationError

	validationErr = request.PathParams.Get("pk", &pk) // 路由传参
	if validationErr != nil {
		return validationErr.Response()
	}

	ctx, cancel := mongodb.WithTimeout(5)
	defer cancel()

	// 执行查询操作
	instance := &DocumentModel{}

	collection := instance.Collection()

	filter := bson.M{"_id": pk}

	err := collection.FindOne(ctx, filter).Decode(instance)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return response.Data{
				Code:    http.StatusNotFound,
				Message: "数据不存在",
				Data:    nil,
			}
		}
		return response.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询失败",
			Data:    err.Error(),
		}
	}

	err = instance.Delete(ctx)
	if err != nil {
		return response.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}

	return response.Data{
		Code:    http.StatusOK,
		Message: "删除成功",
		Data:    nil,
	}
}
