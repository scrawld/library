package ginx

import (
	"reflect"

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
