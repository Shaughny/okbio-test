package utils

import (
	"github.com/go-playground/validator/v10"
	"net"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// ValidateIP checks if a string is a valid IPv4 or IPv6 address
func isValidIP(fl validator.FieldLevel) bool {
	ip := fl.Field().String()
	return net.ParseIP(ip) != nil
}

func NewValidator() *CustomValidator {
	v := validator.New()
	err := v.RegisterValidation("valid_ip", isValidIP)
	if err != nil {
		panic(err)
	}
	return &CustomValidator{validator: v}
}
