package utils

import (
	"reflect"

	"github.com/Chetuks/gym-app/common/models"
	"github.com/go-playground/validator"
)

func ValidateStruct(s interface{}, isWarning bool) *[]models.ValidationError {
	v := validator.New()
	if isWarning {
		v.RegisterValidation("optionalWarn", validateOptionalWarn)
	} else {
		// register a dummy function to avoid custom tag for warning
		v.RegisterValidation("optionalWarn", func(fl validator.FieldLevel) bool { return true })
	}
	var errors []models.ValidationError
	if valErr := v.Struct(s); valErr != nil {
		for _, err := range valErr.(validator.ValidationErrors) {
			var element models.ValidationError
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, element)
		}
		return &errors
	}
	return nil
}

// Custom validation function for optional fields that issue warnings
func validateOptionalWarn(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		return field.String() != ""
	case reflect.Float32, reflect.Float64:
		return field.Float() != 0
	case reflect.Int, reflect.Int64:
		return field.Int() != 0
	default:
		return true
	}
}
