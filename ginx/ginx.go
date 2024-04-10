package ginx

import (
	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/shopspring/decimal"
)

var trans, _ = ut.New(en.New()).GetTranslator("en")

func init() {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		enTranslations.RegisterDefaultTranslations(validate, trans)             // 修改语言
		validate.RegisterCustomTypeFunc(ValidateDecimalType, decimal.Decimal{}) // 将Decimal转成float64进行校验
	}
}

func translateValidationError(err error) error {
	if err == nil {
		return nil
	}
	switch errs := err.(type) {
	case validator.ValidationErrors:
		for _, v := range errs {
			return errors.New(v.Translate(trans))
		}
	case binding.SliceValidationError:
		for _, v := range errs {
			if validationErrors, ok := v.(validator.ValidationErrors); ok {
				for _, ve := range validationErrors {
					return errors.New(ve.Translate(trans))
				}
			}
		}
	}
	return err
}
