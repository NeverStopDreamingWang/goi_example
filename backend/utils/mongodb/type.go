package mongodb

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MongoDB 时间类型
type ISODate time.Time

func (self ISODate) Time() time.Time {
	return time.Time(self)
}

func (self ISODate) MarshalJSON() ([]byte, error) {
	tt := time.Time(self)
	if tt.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + tt.Format(time.DateTime) + `"`), nil
}

func (self *ISODate) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" {
		*self = ISODate(time.Time{})
		return nil
	}

	// 优先尝试 RFC3339 格式
	tt, err := time.Parse(time.RFC3339, s)
	if err == nil {
		*self = ISODate(tt)
		return nil
	}

	// 尝试 DateTime 格式 (2006-01-02 15:04:05)
	tt, err = time.Parse(time.DateTime, s)
	if err == nil {
		*self = ISODate(tt)
		return nil
	}

	return errors.New("时间格式错误，支持格式: 2006-01-02 15:04:05 或 ISO8601")
}

// 实现 bson.ValueMarshaler 接口
func (self ISODate) MarshalBSONValue() (byte, []byte, error) {
	typ, data, err := bson.MarshalValue(time.Time(self))
	return byte(typ), data, err
}

// 实现 bson.ValueUnmarshaler 接口
func (self *ISODate) UnmarshalBSONValue(t byte, data []byte) error {
	var tm time.Time
	err := bson.UnmarshalValue(bson.Type(t), data, &tm)
	if err != nil {
		return err
	}
	*self = ISODate(tm)
	return nil
}

// 实现 goi.AnyUnmarshaler 接口
func (self *ISODate) UnmarshalAny(value any) error {
	if value == nil {
		return nil
	}
	switch typeValue := value.(type) {
	case ISODate:
		*self = typeValue
		return nil
	case time.Time:
		*self = ISODate(typeValue)
		return nil
	case string:
		// 优先尝试 RFC3339 格式
		t, err := time.Parse(time.RFC3339, typeValue)
		if err == nil {
			*self = ISODate(t)
			return nil
		}
		// 尝试 "2006-01-02 15:04:05" 格式
		t, err = time.Parse(time.DateTime, typeValue)
		if err == nil {
			*self = ISODate(t)
			return nil
		}
		return nil
	}
	return nil
}

// ObjectIDList
type ObjectIDList []*bson.ObjectID

func (self *ObjectIDList) UnmarshalAny(value any) error {
	switch v := value.(type) {
	case []any:
		list := make(ObjectIDList, 0, len(v))
		for _, a := range v {
			s, ok := a.(string)
			if !ok {
				continue
			}
			id, err := bson.ObjectIDFromHex(s)
			if err != nil {
				return err
			}
			list = append(list, &id)
		}
		*self = list
		return nil
	}
	return nil
}
