package mongodb

import (
	"context"
	"errors"
	"net/http"

	"goi_example/backend/utils/mongodb"

	"github.com/NeverStopDreamingWang/goi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询标准失败",
			Results: err.Error(),
		}
	}
	defer cursor.Close(ctx)

	var documentList []*DocumentModel
	if err = cursor.All(ctx, &documentList); err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询标准失败",
			Results: err.Error(),
		}
	}

	// 获取总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		total = 0
	}

	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: map[string]any{
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

	document := &DocumentModel{
		Name:    params.Name,
		Content: params.Content,
	}

	// 参数验证
	err := document.Validate()
	if err != nil {
		return goi.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Results: nil,
		}
	}

	// 创建
	err = mongodb.WithTimeoutCtx(5, func(ctx context.Context) error {
		return document.Create(ctx)
	})
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Results: nil,
		}
	}

	return goi.Data{
		Code:    http.StatusOK,
		Message: "创建成功",
		Results: document,
	}
}

func retrieveView(request *goi.Request) any {
	var pk primitive.ObjectID // object_id 转换器将字符串转换为 primitive.ObjectID 类型
	var validationErr goi.ValidationError
	validationErr = request.PathParams.Get("pk", &pk) // 路由转换器自动转换
	if validationErr != nil {
		return validationErr.Response()
	}

	database := mongodb.Database()

	collection := database.Collection("document")

	filter := bson.M{"_id": pk}
	// 执行查询操作
	document := &DocumentModel{}

	ctx, cancel := mongodb.WithTimeout(5)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(document)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return goi.Data{
				Code:    http.StatusNotFound,
				Message: "数据不存在",
				Results: nil,
			}
		}
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询失败",
			Results: err.Error(),
		}
	}

	return goi.Data{
		Code:    http.StatusOK,
		Message: "",
		Results: document,
	}
}

// 参数验证
type updateValidParams struct {
	Name    *string `name:"name" type:"string"`
	Content *string `name:"content" type:"string"`
}

func updateView(request *goi.Request) any {
	var pk primitive.ObjectID
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

	// 执行查询操作
	instance := &DocumentModel{}

	collection := instance.Collection()

	filter := bson.M{"_id": pk}

	ctx, cancel := mongodb.WithTimeout(5)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(instance)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return goi.Data{
				Code:    http.StatusNotFound,
				Message: "数据不存在",
				Results: nil,
			}
		}
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询失败",
			Results: err.Error(),
		}
	}

	document := &DocumentModel{
		Id:      instance.Id,
		Name:    params.Name,
		Content: params.Content,
	}
	err = document.Validate()
	if err != nil {
		return goi.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Results: nil,
		}
	}
	err = instance.Update(ctx, document)
	if err != nil {
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Results: nil,
		}
	}

	return goi.Data{
		Code:    http.StatusOK,
		Message: "修改成功",
		Results: instance,
	}
}

func deleteView(request *goi.Request) any {
	var pk primitive.ObjectID
	var validationErr goi.ValidationError

	validationErr = request.PathParams.Get("pk", &pk) // 路由传参
	if validationErr != nil {
		return validationErr.Response()
	}

	database := mongodb.Database()

	// 获取集合
	collection := database.Collection("document")

	// 执行查询操作
	instance := &DocumentModel{}

	filter := bson.M{"_id": pk}

	ctx, cancel := mongodb.WithTimeout(5)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(instance)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return goi.Data{
				Code:    http.StatusNotFound,
				Message: "数据不存在",
				Results: nil,
			}
		}
		return goi.Data{
			Code:    http.StatusInternalServerError,
			Message: "查询失败",
			Results: err.Error(),
		}
	}

	err = instance.Delete(ctx)
	if err != nil {
		return goi.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Results: nil,
		}
	}

	return goi.Data{
		Code:    http.StatusOK,
		Message: "删除成功",
		Results: nil,
	}
}
