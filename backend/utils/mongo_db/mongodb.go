package mongo_db

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/NeverStopDreamingWang/goi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB 配置
type ConfigModel struct {
	Uri      string `json:"uri"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

var Config *ConfigModel

var MongoDB *mongo.Client

func Connect() {
	if Config == nil {
		panic("MongoDB 配置未初始化")
	}
	var err error
	clientOptions := options.Client()
	clientOptions.ApplyURI(Config.Uri) // MongoDB URI
	clientOptions.SetAuth(options.Credential{
		Username: Config.Username,
		Password: Config.Password,
	})
	clientOptions.SetConnectTimeout(5 * time.Second)
	clientOptions.SetMaxPoolSize(10)

	MongoDB, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		msg := fmt.Sprintf("MongoDB 连接失败: %v", err)
		goi.Log.Error(msg)
		panic(msg)
	}

	// 确保连接成功
	err = MongoDB.Ping(context.Background(), nil)
	if err != nil {
		msg := fmt.Sprintf("MongoDB 连接失败: %v", err)
		goi.Log.Error(msg)
		panic(msg)
	}

}
func Close() error {
	err := MongoDB.Disconnect(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// 获取数据库链接对象
func Database() *mongo.Database {
	return MongoDB.Database(Config.Database)
}

// 事务操作
func WithTransaction(ctx context.Context, transactionFunc func(sessionContext mongo.SessionContext, args ...interface{}) error, args ...interface{}) error {
	var err error
	// 开始事务
	session, err := MongoDB.StartSession()
	if err != nil {
		return errors.New("操作数据库错误")
	}
	defer session.EndSession(ctx)

	// 开始一个会话并定义事务函数
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		// 开始事务
		err = session.StartTransaction()
		if err != nil {
			return errors.New("操作数据库错误")
		}

		// 执行一些数据库操作
		err = transactionFunc(sessionContext, args...)
		if err != nil {
			_ = session.AbortTransaction(sessionContext)
			return err
		}

		// 提交事务
		err = session.CommitTransaction(sessionContext)
		if err != nil {
			_ = session.AbortTransaction(sessionContext)
			return errors.New("操作数据库错误")
		}
		return nil
	})
	// 处理事务错误
	if err != nil {
		return err
	}
	return nil
}

func WithTimeout(second int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(second)*time.Second)
}

func UpdateMap(validated_data interface{}) bson.M {
	update := bson.M{}

	validatedDataValue := reflect.ValueOf(validated_data)
	if validatedDataValue.Kind() == reflect.Ptr {
		validatedDataValue = validatedDataValue.Elem()
	}

	validatedDataType := validatedDataValue.Type()

	for i := 0; i < validatedDataValue.NumField(); i++ {
		field := validatedDataValue.Field(i)
		fieldType := validatedDataType.Field(i)

		if field.IsZero() {
			continue
		}
		if !field.CanSet() {
			continue
		}

		value := field.Interface()

		// 使用 bson tag 作为字段名（没有就用结构体字段名）
		bsonTag := fieldType.Tag.Get("bson")
		if bsonTag == "-" {
			continue
		}
		if bsonTag == "" {
			bsonTag = fieldType.Name
		}

		update[bsonTag] = value
	}
	return update
}
