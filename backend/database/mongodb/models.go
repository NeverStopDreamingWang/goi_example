package mongodb

import (
	"goi_example/backend/utils/mongo_db"

	"github.com/NeverStopDreamingWang/goi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	if mongo_db.Config == nil {
		return
	}
	err := initDocument()
	if err != nil {
		goi.Log.Error(err)
		panic(err)
	}
}

// 模型
type DocumentModel struct {
	Id          *primitive.ObjectID `bson:"id" json:"id"`
	Name        *string             `bson:"name" json:"name"`
	Content     *string             `bson:"content" json:"content"`
	Create_time *string             `bson:"create_Time" json:"create_time"`
	Update_time *string             `bson:"update_Time" json:"update_time"`
}

func initDocument() error {
	initDocumentList := [][]interface{}{
		{"test", "test"},
	}

	database := mongo_db.Database()
	collection := database.Collection("document")

	filter := bson.M{}

	ctx, cancel := mongo_db.WithTimeout(5)
	defer cancel()

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if total > 0 {
		return nil
	}

	for _, item := range initDocumentList {
		var (
			Name    = item[1].(string)
			Content = item[2].(string)
		)
		document := &DocumentModel{
			Name:    &Name,
			Content: &Content,
		}
		// 参数验证
		err = document.Validate()
		if err != nil {
			return err
		}
		// 添加
		err = document.Create()
		if err != nil {
			return err
		}
	}
	return nil
}
