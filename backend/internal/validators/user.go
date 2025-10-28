package validators

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// ValidUsername проверяет username
// Требования: 3-20 символов, только буквы, цифры, дефис, подчеркивание
func ValidUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	if len(username) < 3 || len(username) > 20 {
		return false
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, username)
	return matched
}

// ValidFullName проверяет полное имя
// Требования: 2-50 символов, буквы и пробелы
func ValidFullName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	if len(name) < 2 || len(name) > 50 {
		return false
	}

	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return false
		}
	}

	return true
}
