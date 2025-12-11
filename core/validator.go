package core

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	// Register function to get json tag name for errors
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidationError represents a friendly validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

// ValidateStruct validates a struct and returns friendly errors
func ValidateStruct(s any) []*ValidationError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors []*ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		var element ValidationError
		element.Field = err.Field()
		element.Tag = err.Tag()
		element.Message = msgForTag(err)
		errors = append(errors, &element)
	}
	return errors
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("This field is required")
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters long", fe.Param())
	}
	return fe.Error() // Default error
}
