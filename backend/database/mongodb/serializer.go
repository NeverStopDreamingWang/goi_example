package mongodb

import (
	"errors"
	"time"

	"goi_example/backend/utils"
	"goi_example/backend/utils/mongo_db"

	"github.com/NeverStopDreamingWang/goi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DocumentModel 模型方法
func (self DocumentModel) Validate() error {
	// 自定义验证
	database := mongo_db.Database()
	collection := database.Collection("document")

	filter := bson.M{}
	if self.Id != nil {
		filter["_id"] = bson.M{"$ne": *self.Id}
	}

	if self.Name != nil {
		filter["name"] = self.Name
		ctx, cancel := mongo_db.WithTimeout(5)
		defer cancel()
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return errors.New("查询数据库错误")
		}
		if count > 0 {
			return errors.New("名称重复")
		}
	}
	return nil
}

func (self *DocumentModel) Create() error {
	// 生成新的 ObjectID
	id := primitive.NewObjectID()
	self.Id = &id

	// 设置创建时间和更新时间
	if self.CreateTime == nil {
		CreateTime := goi.GetTime().Format(time.DateTime)
		self.CreateTime = &CreateTime
		self.UpdateTime = &CreateTime
	}

	// 将结构体编码为 BSON 格式
	doc, err := bson.Marshal(self)
	if err != nil {
		return err
	}

	database := mongo_db.Database()
	collection := database.Collection("document")

	// 创建 context 并设置超时时间
	ctx, cancel := mongo_db.WithTimeout(5)
	defer cancel()
	// 使用 InsertOne 插入单个文档
	_, err = collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	return nil
}

func (self *DocumentModel) Update(validated_data *DocumentModel) error {
	updateFields := mongo_db.UpdateMap(validated_data)

	// 如果没有字段需要更新，直接返回
	if len(updateFields) == 0 {
		return nil
	}

	update_time := goi.GetTime().Format(time.DateTime)
	validated_data.UpdateTime = &update_time
	updateFields["update_time"] = validated_data.UpdateTime

	filter := bson.M{"_id": self.Id}

	// 构建更新内容
	update := bson.M{"$set": updateFields}

	database := mongo_db.Database()
	collection := database.Collection("document")

	// 创建 context 并设置超时时间
	ctx, cancel := mongo_db.WithTimeout(5)
	defer cancel()
	// 执行更新操作
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("修改失败")
	}
	utils.Update(self, validated_data)
	return nil
}

func (self DocumentModel) Delete() error {
	if self.Id == nil {
		return errors.New("无效的实例")
	}
	ctx, cancel := mongo_db.WithTimeout(5)
	defer cancel()

	database := mongo_db.Database()
	collection := database.Collection("document")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": *self.Id})
	if err != nil {
		return errors.New("删除失败")
	}
	return nil
}
