package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidCardNumber - проверка номера карты (алгоритм Луна)
func ValidCardNumber(fl validator.FieldLevel) bool {
	cardNumber := fl.Field().String()

	// Убираем пробелы и дефисы
	cardNumber = regexp.MustCompile(`[\s-]`).ReplaceAllString(cardNumber, "")

	// Проверяем что только цифры
	if matched, _ := regexp.MatchString(`^\d+$`, cardNumber); !matched {
		return false
	}

	// Алгоритм Луна
	sum := 0
	double := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	return sum%10 == 0
}

// ValidCVV проверяет CVV код
func ValidCVV(fl validator.FieldLevel) bool {
	cvv := fl.Field().String()
	matched, _ := regexp.MatchString(`^\d{3,4}$`, cvv)
	return matched
}
