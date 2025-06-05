package buyer_service

import (
	"errors"
	"strings"
	"des/backend/internal/models/buyer"
)

func ValidateBuyerRegistration(input *buyer.BuyerRegister) error {
	if strings.TrimSpace(input.EmailAddress) == "" || strings.TrimSpace(input.Password) == "" {
		return errors.New("email and password are required")
	}
	// Add more validation logic as needed (e.g., email regex, password strength)
	return nil
}
