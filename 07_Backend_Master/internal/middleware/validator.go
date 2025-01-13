package middleware

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validation tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidationError represents a validation error
type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ValidateRequest validates the request body against the provided struct
func ValidateRequest(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		if err := validate.Struct(obj); err != nil {
			var errors []ValidationError
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, ValidationError{
					Field: err.Field(),
					Error: getValidationErrorMsg(err),
				})
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"errors": errors,
			})
			return
		}

		c.Next()
	}
}

// getValidationErrorMsg returns a human-readable validation error message
func getValidationErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Should be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("Should be at most %s characters long", err.Param())
	case "oneof":
		return fmt.Sprintf("Should be one of: %s", err.Param())
	default:
		return fmt.Sprintf("Failed validation on %s", err.Tag())
	}
}
