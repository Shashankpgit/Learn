package example

import "github.com/google/uuid"

// ApiRequest_RegistrationRequestExample provides an example for registration request
type ApiRequest_RegistrationRequestExample struct {
	ID      string                     `json:"id" example:"api.actor.create"`
	Ver     string                     `json:"ver" example:"5.0.0"`
	Ts      string                     `json:"ts" example:"2023-01-01T00:00:00Z"`
	Params  ApiRequestParamsExample    `json:"params"`
	Request RegistrationRequestExample `json:"request"`
}

// ApiRequestParamsExample provides example parameters
type ApiRequestParamsExample struct {
	MsgID string `json:"msgid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// RegistrationRequestExample provides an example for registration request payload
type RegistrationRequestExample struct {
	UniversalIdentifier string `json:"universalIdentifier" example:"actor123"`
	Email               string `json:"email" example:"actor@example.com"`
	FirstName           string `json:"firstName" example:"John"`
	LastName            string `json:"lastName" example:"Doe"`
	PhoneNumber         string `json:"phoneNumber" example:"+1234567890"`
	MasterPublicKey     string `json:"masterPublicKey" example:"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...\n-----END PUBLIC KEY-----"`
	EntityType          string `json:"entityType" example:"Individual"`
	Nationality         string `json:"nationality" example:"US"`
	CountryOfResidence  string `json:"countryOfResidence" example:"US"`
}

// ApiRequest_LoginRequestExample provides an example for login request
type ApiRequest_LoginRequestExample struct {
	ID      string                  `json:"id" example:"api.actor.login"`
	Ver     string                  `json:"ver" example:"5.0.0"`
	Ts      string                  `json:"ts" example:"2023-01-01T00:00:00Z"`
	Params  ApiRequestParamsExample `json:"params"`
	Request LoginRequestExample     `json:"request"`
}

// LoginRequestExample provides an example for login request payload
type LoginRequestExample struct {
	Username string `json:"username" example:"actor123"`
	Password string `json:"password" example:"password123"`
}

// ApiRequest_UpdateActorRequestExample provides an example for update request
type ApiRequest_UpdateActorRequestExample struct {
	ID      string                    `json:"id" example:"api.actor.update"`
	Ver     string                    `json:"ver" example:"5.0.0"`
	Ts      string                    `json:"ts" example:"2023-01-01T00:00:00Z"`
	Params  ApiRequestParamsExample   `json:"params"`
	Request UpdateActorRequestExample `json:"request"`
}

// UpdateActorRequestExample provides an example for update request payload
type UpdateActorRequestExample struct {
	FirstName   string `json:"firstName" example:"Jane"`
	LastName    string `json:"lastName" example:"Smith"`
	PhoneNumber string `json:"phoneNumber" example:"+0987654321"`
}

// ApiRequest_ForgotPasswordRequestExample provides an example for forgot password request
type ApiRequest_ForgotPasswordRequestExample struct {
	ID      string                       `json:"id" example:"api.actor.forgotPassword"`
	Ver     string                       `json:"ver" example:"5.0.0"`
	Ts      string                       `json:"ts" example:"2023-01-01T00:00:00Z"`
	Params  ApiRequestParamsExample      `json:"params"`
	Request ForgotPasswordRequestExample `json:"request"`
}

// ForgotPasswordRequestExample provides an example for forgot password request payload
type ForgotPasswordRequestExample struct {
	Email string `json:"email" example:"actor@example.com"`
}

// ApiRequest_ResolveRequestExample provides an example for resolve request
type ApiRequest_ResolveRequestExample struct {
	ID      string                  `json:"id" example:"api.actor.resolve"`
	Ver     string                  `json:"ver" example:"5.0.0"`
	Ts      string                  `json:"ts" example:"2023-01-01T00:00:00Z"`
	Params  ApiRequestParamsExample `json:"params"`
	Request ResolveRequestExample   `json:"request"`
}

// ResolveRequestExample provides an example for resolve request payload
type ResolveRequestExample struct {
	UniversalIdentifier string `json:"universalIdentifier" example:"user123"`
}

// ApiResponse_RegistrationSuccessExample provides an example for registration success response
type ApiResponse_RegistrationSuccessExample struct {
	ID       string                             `json:"id" example:"api.actor.create"`
	Ver      string                             `json:"ver" example:"5.0.0"`
	Ts       string                             `json:"ts" example:"2023-01-01T00:00:00Z"`
	Params   ApiResponseParamsExample           `json:"params"`
	Response RegistrationSuccessResponseExample `json:"response"`
}

// ApiResponseParamsExample provides example response parameters
type ApiResponseParamsExample struct {
	MsgID  string `json:"msgid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Status string `json:"status" example:"successful"`
}

// RegistrationSuccessResponseExample provides an example for registration success response payload
type RegistrationSuccessResponseExample struct {
	DID     uuid.UUID `json:"did" example:"123e4567-e89b-12d3-a456-426614174000"`
	Message string    `json:"message" example:"Actor registered successfully"`
}

// ApiResponse_AuthResponseExample provides an example for auth response
type ApiResponse_AuthResponseExample struct {
	ID       string                   `json:"id" example:"api.actor.login"`
	Ver      string                   `json:"ver" example:"5.0.0"`
	Ts       string                   `json:"ts" example:"2023-01-01T00:00:00Z"`
	Params   ApiResponseParamsExample `json:"params"`
	Response AuthResponseExample      `json:"response"`
}

// AuthResponseExample provides an example for auth response payload
type AuthResponseExample struct {
	AccessToken string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"tokenType" example:"Bearer"`
	ExpiresIn   int    `json:"expiresIn" example:"3600"`
}

// ApiResponse_ActorProfileExample provides an example for actor profile response
type ApiResponse_ActorProfileExample struct {
	ID       string                   `json:"id" example:"api.actor.getProfile"`
	Ver      string                   `json:"ver" example:"5.0.0"`
	Ts       string                   `json:"ts" example:"2023-01-01T00:00:00Z"`
	Params   ApiResponseParamsExample `json:"params"`
	Response ActorProfileExample      `json:"response"`
}

// ActorProfileExample provides an example for actor profile response payload
type ActorProfileExample struct {
	ID                  uuid.UUID `json:"id" example:"e088d183-9eea-4a11-8d5d-74d7ec91bdf5"`
	DID                 uuid.UUID `json:"did" example:"123e4567-e89b-12d3-a456-426614174000"`
	UniversalIdentifier string    `json:"universalIdentifier" example:"user123"`
	Email               string    `json:"email" example:"actor@example.com"`
	FirstName           string    `json:"firstName" example:"John"`
	LastName            string    `json:"lastName" example:"Doe"`
	PhoneNumber         *string   `json:"phoneNumber" example:"+1234567890"`
	MasterPublicKey     string    `json:"masterPublicKey" example:"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...\n-----END PUBLIC KEY-----"`
	VerificationLevel   string    `json:"verificationLevel" example:"Tier0_Unverified"`
	EntityType          string    `json:"entityType" example:"Individual"`
	Nationality         *string   `json:"nationality" example:"US"`
	CountryOfResidence  *string   `json:"countryOfResidence" example:"US"`
}

// ApiResponse_ResolveResponseExample provides an example for resolve response
type ApiResponse_ResolveResponseExample struct {
	ID       string                   `json:"id" example:"api.actor.resolve"`
	Ver      string                   `json:"ver" example:"5.0.0"`
	Ts       string                   `json:"ts" example:"2023-01-01T00:00:00Z"`
	Params   ApiResponseParamsExample `json:"params"`
	Response ResolveResponseExample   `json:"response"`
}

// ResolveResponseExample provides an example for resolve response payload
type ResolveResponseExample struct {
	UniversalIdentifier string    `json:"universalIdentifier" example:"user123"`
	MasterPublicKey     string    `json:"masterPublicKey" example:"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...\n-----END PUBLIC KEY-----"`
	DID                 uuid.UUID `json:"did" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// ApiResponse_SuccessExample provides an example for success response
type ApiResponse_SuccessExample struct {
	ID       string                   `json:"id" example:"api.actor.signout"`
	Ver      string                   `json:"ver" example:"5.0.0"`
	Ts       string                   `json:"ts" example:"2023-01-01T00:00:00Z"`
	Params   ApiResponseParamsExample `json:"params"`
	Response SuccessResponseExample   `json:"response"`
}

// SuccessResponseExample provides an example for success response payload
type SuccessResponseExample struct {
	Message string `json:"message" example:"Signout successful"`
}
