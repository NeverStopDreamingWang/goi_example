package goi_example

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/NeverStopDreamingWang/goi/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	_ goi.Validator = phoneValidator{}
	_ goi.Validator = objectIdValidator{}
)

// 使用显式文件列表运行，避免加载项目配置和数据库初始化：
// go test ./backend/goi_example/validator.go ./backend/goi_example/validator_test.go
func TestPhoneValidatorParse(t *testing.T) {
	phone, err := strconv.Atoi("13800000000")
	validInputs := []any{"13800000000"}
	if err == nil {
		validInputs = append(validInputs, phone)
	}
	validator, ok := goi.GetValidator("phone")
	if !ok {
		t.Fatal("phone validator not registered")
	}
	if _, ok := validator.(phoneValidator); !ok {
		t.Fatalf("registered phone validator = %T; want phoneValidator", validator)
	}
	for _, input := range validInputs {
		got, validationErr := validator.Parse(input)
		value, isString := got.(string)
		if validationErr != nil || !isString || value != "13800000000" {
			t.Fatalf("Parse(%#v) = %#v, %v; want full phone string", input, got, validationErr)
		}
	}
	for _, input := range []any{"invalid", "1380000000", "13800000000/extra", true, int64(13800000000), nil} {
		got, validationErr := validator.Parse(input)
		if got != nil || validationErr == nil || validationErr.Response().Status != http.StatusBadRequest {
			t.Fatalf("Parse(%#v) = %#v, %v; want nil and status 400", input, got, validationErr)
		}
	}
}

func TestObjectIDValidatorParse(t *testing.T) {
	const hex = "507f1f77bcf86cd799439011"
	want, err := bson.ObjectIDFromHex(hex)
	if err != nil {
		t.Fatal(err)
	}
	validator, ok := goi.GetValidator("object_id")
	if !ok {
		t.Fatal("object_id validator not registered")
	}
	if _, ok := validator.(objectIdValidator); !ok {
		t.Fatalf("registered object_id validator = %T; want objectIdValidator", validator)
	}
	for _, input := range []string{hex, strings.ToUpper(hex)} {
		got, validationErr := validator.Parse(input)
		value, isObjectID := got.(bson.ObjectID)
		if validationErr != nil || !isObjectID || value != want {
			t.Fatalf("Parse(%q) = %#v, %v; want bson.ObjectID %s", input, got, validationErr, want.Hex())
		}
	}
	for _, input := range []any{hex[:23], hex + "0", strings.Repeat("g", 24), " " + hex, 123, want, nil} {
		got, validationErr := validator.Parse(input)
		if got != nil || validationErr == nil || validationErr.Response().Status != http.StatusBadRequest {
			t.Fatalf("Parse(%#v) = %#v, %v; want nil and status 400", input, got, validationErr)
		}
	}
}

func TestCustomValidatorsParseParams(t *testing.T) {
	const phone = "13800000000"
	const hex = "507f1f77bcf86cd799439011"
	wantID, err := bson.ObjectIDFromHex(hex)
	if err != nil {
		t.Fatal(err)
	}
	var params struct {
		Phone string        `name:"phone" type:"phone" required:"true"`
		ID    bson.ObjectID `name:"id" type:"object_id" required:"true"`
	}
	values := goi.Params{"phone": phone, "id": strings.ToUpper(hex)}
	if validationErr := values.ParseParams(&params); validationErr != nil {
		t.Fatal(validationErr)
	}
	if params.Phone != phone || params.ID != wantID {
		t.Fatalf("parsed params = %+v; want phone %s and ID %s", params, phone, wantID.Hex())
	}
}
