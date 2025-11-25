package utils

import (
	"reflect"
)

func Update(instance any, validated_data any) {
	instanceValue := reflect.ValueOf(instance)
	validatedDataValue := reflect.ValueOf(validated_data)

	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}
	if validatedDataValue.Kind() == reflect.Ptr {
		validatedDataValue = validatedDataValue.Elem()
	}
	instanceType := instanceValue.Type()

	for i := 0; i < instanceType.NumField(); i++ {
		field := instanceType.Field(i)
		instanceField := instanceValue.Field(i)
		validatedField := validatedDataValue.FieldByName(field.Name)

		if validatedField.Kind() == reflect.Ptr && validatedField.IsNil() {
			continue
		}
		if !instanceField.CanSet() {
			continue
		}
		instanceField.Set(validatedField)
	}
}
