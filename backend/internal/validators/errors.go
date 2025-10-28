package validators

import (
	"fmt"
	"strings"

	"github.com/ZAUakaAlexey/backend_go/internal/responses"
	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(err error) []responses.ValidationError {
	var errors []responses.ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, responses.ValidationError{
				Field:   toSnakeCase(e.Field()),
				Message: getErrorMessage(e),
				Tag:     e.Tag(),
			})
		}
	}

	return errors
}

func getErrorMessage(e validator.FieldError) string {
	field := e.Field()

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "strongpassword":
		return "Password must be at least 8 characters and contain at least 3 of: uppercase, lowercase, number, special character"
	case "username":
		return "Username must be 3-20 characters and contain only letters, numbers, dash and underscore"
	case "fullname":
		return "Full name must be 2-50 characters and contain only letters and spaces"
	case "phone":
		return "Invalid phone number format"
	case "notempty":
		return fmt.Sprintf("%s cannot be empty", field)
	case "alphaspace":
		return fmt.Sprintf("%s can only contain letters and spaces", field)
	case "cardnumber":
		return "Invalid card number"
	case "cvv":
		return "CVV must be 3 or 4 digits"
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder

	for i, char := range s {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(char)
	}

	return strings.ToLower(result.String())
}
