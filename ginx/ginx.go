package ginx

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/scrawld/library/types"
	"github.com/shopspring/decimal"
)

func init() {
	binding.Validator = &defaultValidator{}
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		validate.RegisterCustomTypeFunc(ValidateDecimalType, decimal.Decimal{})             // 将Decimal转成float64进行校验
		validate.RegisterCustomTypeFunc(ValidateCustomTimeType, types.Date{}, types.Time{}) // 将自定义Time转成time.Time进行校验
		validate.RegisterValidation("trimspace", ValidateTrimSpace)                         // 去除字符串两端空格
	}
}
