package ginx

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func init() {
	binding.Validator = &defaultValidator{}
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		validate.RegisterCustomTypeFunc(ValidateDecimalType, decimal.Decimal{}) // 将Decimal转成float64进行校验
		validate.RegisterValidation("trimspace", ValidateTrimSpace)             // 去除字符串两端空格
	}
}
