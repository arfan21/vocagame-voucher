package validation

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni        *ut.UniversalTranslator
	validate   *validator.Validate
	translator ut.Translator

	uniOnce        sync.Once
	validateOnce   sync.Once
	translatorOnce sync.Once
)

type ErrValidationMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ErrValidationMessage) Error() string {
	return fmt.Sprintf("%v: %v", e.Field, e.Message)
}

type ErrsValidation []ErrValidationMessage

func (e ErrsValidation) Error() string {
	var messages []string
	for _, err := range e {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, ", ")
}

func Validate[T any](modelValidate T) error {
	uniOnce.Do(func() {
		en := en.New()
		uni = ut.New(en, en)
	})

	validateOnce.Do(func() {
		validate = validator.New()
	})

	translatorOnce.Do(func() {
		translatorUni, _ := uni.GetTranslator("en")
		translator = translatorUni
		en_translations.RegisterDefaultTranslations(validate, translator)
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := validate.Struct(modelValidate)
	if err != nil {
		var messages ErrsValidation

		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()

			messages = append(messages, ErrValidationMessage{
				Field:   fieldName,
				Message: err.Translate(translator),
			})
		}

		return messages
	}

	return nil
}
