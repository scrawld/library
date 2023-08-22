package ginx

import (
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

// ValidateDecimalType 将Decimal转成float64进行校验
func ValidateDecimalType(field reflect.Value) interface{} {
	fieldDecimal, ok := field.Interface().(decimal.Decimal)
	if ok {
		value, _ := fieldDecimal.Float64()
		return value
	}
	return field
}

// ValidateCustomTimeType 将自定义Time转成time.Time进行校验
func ValidateCustomTimeType(field reflect.Value) interface{} {
	type CustomTime interface{ ToTime() time.Time }

	fieldTime, ok := field.Interface().(CustomTime)
	if ok {
		value := fieldTime.ToTime()
		return value
	}
	return field.Interface()
}

// ValidateTrimSpace 去除string两端空格
func ValidateTrimSpace(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.String {
		if !field.CanSet() {
			return false // 需要传递指针
		}
		field.SetString(strings.TrimSpace(field.String()))
	}
	return true
}
