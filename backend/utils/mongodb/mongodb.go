package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/NeverStopDreamingWang/goi/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	MongoDB, err = mongo.Connect(clientOptions)
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
func WithTransaction(ctx context.Context, transactionFunc func(sessionContext context.Context, args ...any) error, args ...any) error {
	session, err := MongoDB.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessionContext context.Context) (any, error) {
		return nil, transactionFunc(sessionContext, args...)
	})
	return err
}

func WithTimeout(second int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(second)*time.Second)
}

func WithTimeoutCtx(second int, fn func(ctx context.Context) error) error {
	ctx, cancel := WithTimeout(second)
	defer cancel()
	return fn(ctx)
}

// UpdateMap 将结构体转换为 bson.M，跳过零值字段
//
// 参数:
//   - validated_data any: 要转换的结构体实例
//   - allowZeroFields ...string: 允许为零值的字段名列表
//
// 返回:
//   - bson.M: 转换后的 bson.M
func UpdateMap(validated_data any, allowZeroFields ...string) bson.M {
	update := bson.M{}

	allowZero := make(map[string]struct{})
	for _, f := range allowZeroFields {
		allowZero[f] = struct{}{}
	}

	validatedDataValue := reflect.ValueOf(validated_data)
	if validatedDataValue.Kind() == reflect.Ptr {
		validatedDataValue = validatedDataValue.Elem()
	}

	validatedDataType := validatedDataValue.Type()

	for i := 0; i < validatedDataValue.NumField(); i++ {
		field := validatedDataValue.Field(i)

		if !field.CanInterface() {
			continue
		}

		fieldType := validatedDataType.Field(i)

		bsonName := getBsonName(fieldType)
		if bsonName == "" {
			continue
		}

		_, isAllow := allowZero[bsonName]
		if field.IsZero() && !isAllow {
			continue
		}
		update[bsonName] = field.Interface()
	}
	return update
}

func getBsonName(fieldType reflect.StructField) string {
	tag := fieldType.Tag.Get("bson")
	if tag == "-" {
		return ""
	}
	if tag == "" {
		return fieldType.Name
	}
	return strings.Split(tag, ",")[0]
}
