package validators

import (
	"context"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const (
	emailRegexStr    = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	usernameRegexStr = `^[a-zA-Z0-9._%+\-]+$`
)

// validateEmailOrUsername validates username and email validation for login request
func (v *Validator) validateEmailOrUsername(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if strings.Contains(value, "@") {
		emailRegex := regexp.MustCompile(emailRegexStr)
		return emailRegex.MatchString(value)
	}
	usernameRegex := regexp.MustCompile(usernameRegexStr)
	return usernameRegex.MatchString(value)
}

// validateUsername validates username validation for signup request
func (v *Validator) validateUsername(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	usernameRegex := regexp.MustCompile(usernameRegexStr)
	return usernameRegex.MatchString(value)
}

// validateStrongPassword validates password validations in signup request
func (v *Validator) validateStrongPassword(fl validator.FieldLevel) bool {
	return v.hasSpecialCharacter(fl) &&
		v.hasNumber(fl) &&
		v.hasUpperCase(fl) &&
		v.hasLowerCase(fl)
}

func (v *Validator) hasSpecialCharacter(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	requiredSpecialCharCount := 1
	var specialCharCount int
	specialChars := "!@#$%^&*()_+-=[]{}|;:',.<>?/"
	for _, char := range password {
		if strings.ContainsRune(specialChars, char) {
			specialCharCount++
		}
	}
	return specialCharCount >= requiredSpecialCharCount
}

func (v *Validator) hasNumber(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	for _, char := range password {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func (v *Validator) hasUpperCase(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var hasUpperCase bool
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
			break
		}
	}
	return hasUpperCase
}

func (v *Validator) hasLowerCase(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var hasLowerCase bool
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowerCase = true
			break
		}
	}
	return hasLowerCase
}

// Uniqueness of Email and Username validation in sign up request
func (v *Validator) validateUniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	user, _ := v.userRepository.FindByEmailOrUsername(context.Background(), email)
	return user == nil
}

func (v *Validator) validateUniqueUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	user, _ := v.userRepository.FindByEmailOrUsername(context.Background(), username)
	return user == nil
}
