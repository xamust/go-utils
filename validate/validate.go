package validate

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	val "github.com/go-playground/validator/v10"
	rutranslations "github.com/go-playground/validator/v10/translations/ru"
)

var (
	validator = val.New(val.WithRequiredStructEnabled())
	trans     ut.Translator
)

func init() {
	rus := ru.New()
	uni := ut.New(rus, rus)

	trans, _ = uni.GetTranslator("ru")

	_ = rutranslations.RegisterDefaultTranslations(validator, trans)
}

type Errors struct {
	errs []error
}

func (e *Errors) Add(err val.FieldError) {
	e.errs = append(e.errs, errors.New(err.Translate(trans)))
}

func (e *Errors) Error() string {
	var msg []string
	for index := range e.errs {
		msg = append(msg, e.errs[index].Error())
	}

	return strings.Join(msg, "; ")
}

func Validate(src any) error {
	return validator.Struct(src)
}
