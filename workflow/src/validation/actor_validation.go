package validation

import (
	"app/src/constants"
	"errors"
)

// RegistrationRequest represents the request payload for actor registration
type RegistrationRequest struct {
	UniversalIdentifier string `json:"universalIdentifier" validate:"required" example:"actor123"`
	Email               string `json:"email" validate:"required,email" example:"actor@example.com"`
	Password            string `json:"password" validate:"required,min=8,password" example:"password123"`
	FirstName           string `json:"firstName" validate:"required" example:"John"`
	LastName            string `json:"lastName" validate:"required" example:"Doe"`
	PhoneNumber         string `json:"phoneNumber,omitempty" validate:"omitempty,e164" example:"+1234567890"`
	MasterPublicKey     string `json:"masterPublicKey" validate:"required" example:"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...\n-----END PUBLIC KEY-----"`
	EntityType          string `json:"entityType" validate:"required,oneof=Individual Business" example:"Individual"`
	// Individual fields
	Nationality        string `json:"nationality,omitempty" validate:"omitempty,len=2" example:"US"`
	CountryOfResidence string `json:"countryOfResidence,omitempty" validate:"omitempty,len=2" example:"US"`
	// Business fields
	CountryOfIncorporation string `json:"countryOfIncorporation,omitempty" validate:"omitempty,len=2" example:"US"`
}

// ValidateEntityTypeFields validates entity type specific required fields
func (r *RegistrationRequest) ValidateEntityTypeFields() error {
	switch r.EntityType {
	case constants.EntityTypeIndividual:
		if r.Nationality == "" {
			return errors.New(constants.ErrNationalityRequiredForIndividual)
		}
		if r.CountryOfResidence == "" {
			return errors.New(constants.ErrCountryOfResidenceRequiredForIndividual)
		}
	case constants.EntityTypeBusiness:
		if r.CountryOfIncorporation == "" {
			return errors.New(constants.ErrCountryOfIncorporationRequiredForBusiness)
		}
	default:
		return errors.New(constants.ErrInvalidEntityType)
	}
	return nil
}

// LoginRequest represents the request payload for actor login
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"alice@finternet"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// UpdateActorRequest represents the request payload for updating actor profile
type UpdateActorRequest struct {
	FirstName   string  `json:"firstName,omitempty" validate:"omitempty" example:"Jane"`
	LastName    string  `json:"lastName,omitempty" validate:"omitempty" example:"Smith"`
	PhoneNumber *string `json:"phoneNumber,omitempty" validate:"omitempty,e164" example:"+0987654321"`
}

// ForgotPasswordRequest represents the request payload for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"actor@example.com"`
}

// ResolveRequest represents the request payload for resolving universal identifier
type ResolveRequest struct {
	UniversalIdentifier string `json:"universalIdentifier" validate:"required" example:"actor123"`
}

// ApiRequest wrapper for all requests
type ApiRequest struct {
	ID      string        `json:"id" validate:"required" example:"api.actor.create"`
	Ver     string        `json:"ver" validate:"required" example:"5.0.0"`
	Ts      string        `json:"ts" validate:"required,datetime" example:"2023-01-01T00:00:00Z"`
	Params  RequestParams `json:"params"`
	Request interface{}   `json:"request"`
}

type RequestParams struct {
	MsgID *string `json:"msgid,omitempty" example:"2f7c3a4e-8bcb-4b44-8cf6-84cc8f1d2a87"`
}

// ApiRequest_RegistrationRequest wraps RegistrationRequest with ApiRequest
type ApiRequest_RegistrationRequest struct {
	ApiRequest
	Request RegistrationRequest `json:"request"`
}

// ApiRequest_LoginRequest wraps LoginRequest with ApiRequest
type ApiRequest_LoginRequest struct {
	ApiRequest
	Request LoginRequest `json:"request"`
}

// ApiRequest_UpdateActorRequest wraps UpdateActorRequest with ApiRequest
type ApiRequest_UpdateActorRequest struct {
	ApiRequest
	Request UpdateActorRequest `json:"request"`
}

// ApiRequest_ForgotPasswordRequest wraps ForgotPasswordRequest with ApiRequest
type ApiRequest_ForgotPasswordRequest struct {
	ApiRequest
	Request ForgotPasswordRequest `json:"request"`
}

// ApiRequest_ResolveRequest wraps ResolveRequest with ApiRequest
type ApiRequest_ResolveRequest struct {
	ApiRequest
	Request ResolveRequest `json:"request"`
}

// ApiRequest_Empty represents an empty request body
type ApiRequest_Empty struct {
	ApiRequest
	Request struct{} `json:"request"`
}
