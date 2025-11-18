package response

import "github.com/google/uuid"

// ApiResponse represents the base response structure
type ApiResponse struct {
	ID       string            `json:"id" example:"api.actor.create"`
	Ver      string            `json:"ver" example:"v1"`
	Ts       string            `json:"ts" example:"2023-01-01T00:00:00Z"`
	Params   ApiResponseParams `json:"params"`
	Response interface{}       `json:"response"`
}

// ApiResponseParams represents the parameters in API responses
type ApiResponseParams struct {
	MsgID  string `json:"msgid,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
	Status string `json:"status" example:"successful"`
	Err    string `json:"err,omitempty"`
	ErrMsg string `json:"errmsg,omitempty"`
}

// RegistrationSuccessResponse represents the response for successful actor registration
type RegistrationSuccessResponse struct {
	DID     uuid.UUID `json:"did" example:"123e4567-e89b-12d3-a456-426614174000"`
	Message string    `json:"message" example:"Actor registered successfully"`
}

// AuthResponse represents the response for successful login
type AuthResponse struct {
	AccessToken string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"tokenType" example:"Bearer"`
	ExpiresIn   int    `json:"expiresIn" example:"3600"`
}

// ActorProfile represents the actor profile information
type ActorProfile struct {
	DID                 uuid.UUID `json:"did" example:"123e4567-e89b-12d3-a456-426614174000"`
	UniversalIdentifier string    `json:"universalIdentifier" example:"actor123"`
	Email               string    `json:"email" example:"actor@example.com"`
	FirstName           string    `json:"firstName" example:"John"`
	LastName            string    `json:"lastName" example:"Doe"`
	PhoneNumber         *string   `json:"phoneNumber,omitempty" example:"+1234567890"`
	MasterPublicKey     string    `json:"masterPublicKey" example:"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...\n-----END PUBLIC KEY-----"`
	VerificationLevel   string    `json:"verificationLevel" example:"Tier0_Unverified"`
	EntityType          string    `json:"entityType" example:"Individual"`
	// Individual fields
	Nationality        *string `json:"nationality,omitempty" example:"US"`
	CountryOfResidence *string `json:"countryOfResidence,omitempty" example:"US"`
	// Business fields
	CountryOfIncorporation *string `json:"countryOfIncorporation,omitempty" example:"US"`
}

// ResolveResponse represents the response for resolving universal identifier
type ResolveResponse struct {
	UniversalIdentifier string    `json:"universalIdentifier" example:"user123"`
	MasterPublicKey     string    `json:"masterPublicKey" example:"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...\n-----END PUBLIC KEY-----"`
	DID                 uuid.UUID `json:"did" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// ApiResponse_RegistrationSuccess wraps RegistrationSuccessResponse with ApiResponse
type ApiResponse_RegistrationSuccess struct {
	ApiResponse
	Response RegistrationSuccessResponse `json:"response"`
}

// ApiResponse_AuthResponse wraps AuthResponse with ApiResponse
type ApiResponse_AuthResponse struct {
	ApiResponse
	Response AuthResponse `json:"response"`
}

// ApiResponse_ActorProfile wraps ActorProfile with ApiResponse
type ApiResponse_ActorProfile struct {
	ApiResponse
	Response ActorProfile `json:"response"`
}

// ApiResponse_ResolveResponse wraps ResolveResponse with ApiResponse
type ApiResponse_ResolveResponse struct {
	ApiResponse
	Response ResolveResponse `json:"response"`
}

// ApiResponse_Success wraps SuccessResponse with ApiResponse
type ApiResponse_Success struct {
	ApiResponse
	Response SuccessResponse `json:"response"`
}

// ApiResponse_Error wraps error response with ApiResponse
type ApiResponse_Error struct {
	ApiResponse
	Response interface{} `json:"response"`
}
