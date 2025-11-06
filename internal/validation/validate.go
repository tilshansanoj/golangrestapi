package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var validate = validator.New()

func Validate(i interface{}) error {
	err := validate.Struct(i)
	if err != nil {
		var errMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			errMessages = append(errMessages, fmt.Sprintf("Field '%s' failed validation rule '%s'", err.Field(), err.Tag()))
		}
		return fmt.Errorf("%s", strings.Join(errMessages, "; "))
	}
	return nil
}