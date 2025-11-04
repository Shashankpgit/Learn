package validation

import (
	"app/src/constants"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var customMessages = map[string]string{
	constants.ValidationTagNameRequired: constants.ValidationMsgRequired,
	constants.ValidationTagNameEmail:    constants.ValidationMsgEmail,
	constants.ValidationTagNameMin:      constants.ValidationMsgMin,
	constants.ValidationTagNameMax:      constants.ValidationMsgMax,
	constants.ValidationTagNameLen:      constants.ValidationMsgLen,
	constants.ValidationTagNameNumber:   constants.ValidationMsgNumber,
	constants.ValidationTagNamePositive: constants.ValidationMsgPositive,
	constants.ValidationTagNameAlphanum: constants.ValidationMsgAlphanum,
	constants.ValidationTagNameOneOf:    constants.ValidationMsgOneOf,
	constants.ValidationTagNamePassword: constants.ValidationMsgPassword,
	constants.ValidationTagNameE164:     constants.ValidationMsgE164,
}

func CustomErrorMessages(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return generateErrorMessages(validationErrors)
	}
	return nil
}

func generateErrorMessages(validationErrors validator.ValidationErrors) map[string]string {
	errorsMap := make(map[string]string)
	for _, err := range validationErrors {
		fieldName := err.StructNamespace()
		tag := err.Tag()

		customMessage := customMessages[tag]
		if customMessage != "" {
			errorsMap[fieldName] = formatErrorMessage(customMessage, err, tag)
		} else {
			errorsMap[fieldName] = defaultErrorMessage(err)
		}
	}
	return errorsMap
}

func formatErrorMessage(customMessage string, err validator.FieldError, tag string) string {
	if tag == constants.ValidationTagNameMin || tag == constants.ValidationTagNameMax || tag == constants.ValidationTagNameLen {
		return fmt.Sprintf(customMessage, err.Field(), err.Param())
	}
	if tag == constants.ValidationTagNamePassword {
		return fmt.Sprintf(customMessage, err.Field(), constants.PasswordMinLength)
	}
	return fmt.Sprintf(customMessage, err.Field())
}

func defaultErrorMessage(err validator.FieldError) string {
	return fmt.Sprintf(constants.ValidationMsgDefault, err.Field(), err.Tag())
}

// NewValidator creates a new validator instance with custom validations
// This is a constructor function for dependency injection
func NewValidator() *validator.Validate {
	validate := validator.New()

	if err := validate.RegisterValidation(constants.ValidationTagNamePassword, Password); err != nil {
		return nil
	}

	if err := validate.RegisterValidation(constants.ValidationTagNameE164, E164); err != nil {
		return nil
	}

	return validate
}
