package validation

import (
	"app/src/constants"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Password(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	if ok {
		hasDigit := regexp.MustCompile(constants.RegexDigit).MatchString(value)
		hasLetter := regexp.MustCompile(constants.RegexLetter).MatchString(value)

		if !hasDigit || !hasLetter {
			return false
		}
	}

	return true
}

func E164(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	if ok {
		// E.164 format: + followed by country code and subscriber number
		// Basic validation: starts with +, followed by digits, length between min-max
		if strings.HasPrefix(value, constants.PhonePrefixPlus) &&
			len(value) >= constants.PhoneNumberMinLength &&
			len(value) <= constants.PhoneNumberMaxLength {
			digitsOnly := strings.ReplaceAll(value[1:], constants.PhoneSpace, "")
			hasOnlyDigits := regexp.MustCompile(constants.RegexDigitsOnly).MatchString(digitsOnly)
			return hasOnlyDigits
		}
	}

	return false
}
