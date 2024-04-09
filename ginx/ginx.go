package ginx

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/scrawld/library/types"
	"github.com/shopspring/decimal"
)

var trans, _ = ut.New(en.New()).GetTranslator("en")

func init() {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		enTranslations.RegisterDefaultTranslations(validate, trans)                         // 修改语言
		validate.RegisterCustomTypeFunc(ValidateDecimalType, decimal.Decimal{})             // 将Decimal转成float64进行校验
		validate.RegisterCustomTypeFunc(ValidateCustomTimeType, types.Date{}, types.Time{}) // 将自定义Time转成time.Time进行校验
	}
}
