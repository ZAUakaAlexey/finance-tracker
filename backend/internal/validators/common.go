package validators

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// ValidPhone проверяет номер телефона
func ValidPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Простая проверка, можно расширить для разных форматов
	matched, _ := regexp.MatchString(`^\+?[1-9]\d{1,14}$`, phone)
	return matched
}

// NotEmpty проверяет что строка не пустая после trim
func NotEmpty(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return strings.TrimSpace(value) != ""
}

// AlphaSpace - только буквы и пробелы
func AlphaSpace(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	for _, char := range value {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) {
			return false
		}
	}

	return true
}
