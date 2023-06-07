package validators

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	ptBR_translations "github.com/go-playground/validator/v10/translations/pt_BR"
)

var validate *validator.Validate
var trans ut.Translator

func NewValidator(language string) {
	var uni *ut.UniversalTranslator
	validate = validator.New()

	switch language {
	case "en":
		uni = ut.New(en.New())
		trans, _ = uni.GetTranslator(language)
		en_translations.RegisterDefaultTranslations(validate, trans)
	case "pt_BR":
		uni = ut.New(pt_BR.New())
		trans, _ = uni.GetTranslator(language)
		ptBR_translations.RegisterDefaultTranslations(validate, trans)
	default:
		uni = ut.New(en.New())
		trans, _ = uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(validate, trans)
	}
}

func Validate(data interface{}) interface{} {
	err := validate.Struct(data)
	if err != nil {
		errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.Translate(trans))
		}

		return errors
	}

	return nil
}
