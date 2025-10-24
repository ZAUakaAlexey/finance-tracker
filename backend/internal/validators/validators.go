package validators

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Passwords
		if err := v.RegisterValidation("strongpassword", StrongPassword); err != nil {
			return fmt.Errorf("failed to register strongpassword: %w", err)
		}

		// User credentials
		if err := v.RegisterValidation("username", ValidUsername); err != nil {
			return fmt.Errorf("failed to register username: %w", err)
		}

		if err := v.RegisterValidation("fullname", ValidFullName); err != nil {
			return fmt.Errorf("failed to register fullname: %w", err)
		}

		// Common
		if err := v.RegisterValidation("phone", ValidPhone); err != nil {
			return fmt.Errorf("failed to register phone: %w", err)
		}

		if err := v.RegisterValidation("notempty", NotEmpty); err != nil {
			return fmt.Errorf("failed to register notempty: %w", err)
		}

		if err := v.RegisterValidation("alphaspace", AlphaSpace); err != nil {
			return fmt.Errorf("failed to register alphaspace: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to get validator engine")
}
