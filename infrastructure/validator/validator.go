package validator

import (
	"fmt"
	"reflect"
	"strings"

	"gin-boilerplate/infrastructure/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// Register custom tag name function to use json tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct validates a struct and returns formatted error
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}
		return fmt.Errorf(strings.Join(validationErrors, "; "))
	}
	return nil
}

// BindAndValidate binds request data and validates it
func BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.NewBadRequestError("Invalid JSON format: " + err.Error())
	}

	if err := ValidateStruct(obj); err != nil {
		return errors.NewValidationError(err.Error())
	}

	return nil
}

// BindQueryAndValidate binds query parameters and validates them
func BindQueryAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.NewBadRequestError("Invalid query parameters: " + err.Error())
	}

	if err := ValidateStruct(obj); err != nil {
		return errors.NewValidationError(err.Error())
	}

	return nil
}

// BindURIAndValidate binds URI parameters and validates them
func BindURIAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindUri(obj); err != nil {
		return errors.NewBadRequestError("Invalid URI parameters: " + err.Error())
	}

	if err := ValidateStruct(obj); err != nil {
		return errors.NewValidationError(err.Error())
	}

	return nil
}

// formatValidationError formats validation error message
func formatValidationError(err validator.FieldError) string {
	field := err.Field()
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
