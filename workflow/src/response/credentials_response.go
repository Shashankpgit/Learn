package response

// BaseRequest represents the base structure for API requests
type BaseRequest struct {
	ID     string        `json:"id" example:"api_id"`
	Ver    string        `json:"ver" example:"v1"`
	Ts     string        `json:"ts" example:"2025-09-08T12:00:00Z"`
	Params RequestParams `json:"params"`
}

// BaseResponse represents the base structure for API responses
type BaseResponse struct {
	ID     string         `json:"id" example:"api_id"`
	Ver    string         `json:"ver" example:"v1"`
	Ts     string         `json:"ts" example:"2025-09-08T12:00:00Z"`
	Params ResponseParams `json:"params"`
}

// ResponseParams represents response parameters in the new API format
type ResponseParams struct {
	MsgID    *string `json:"msgid,omitempty"`
	Status   *string `json:"status,omitempty"`
	Err      *string `json:"err,omitempty"`
	ErrMsg   *string `json:"errmsg,omitempty"`
	ResMsgID *string `json:"resmsgid,omitempty"`
}

// RequestParams represents request parameters in the new API format
type RequestParams struct {
	MsgID *string `json:"msgid,omitempty" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// Request is a generic request wrapper
type Request[T any] struct {
	BaseRequest
	Request T `json:"request"`
}

// Response is a generic response wrapper
type Response[T any] struct {
	BaseResponse
	ResponseCode string `json:"responseCode"`
	Response     T      `json:"response"`
}

// AddCredentialSuccessResponse represents the payload returned on successful
// AddCredential operation.
type AddCredentialSuccessResponse struct {
	CredentialID string `json:"credentialId" example:"123e4567-e89b-12d3-a456-426614174000"`
	Status       string `json:"status" example:"Pending"`
	Message      string `json:"message" example:"Credential submitted successfully. Verification is in progress."`
}

// CredentialsSuccessResponse represents a single credential in responses
type CredentialsSuccessResponse struct {
	CredentialID string `json:"credentialId" example:"123e4567-e89b-12d3-a456-426614174000"`
	Type         string `json:"type" example:"VerifiableCredential"`
	Status       string `json:"status" example:"Pending"`
	SubmittedAt  string `json:"submittedAt" example:"2025-10-23T06:25:25.191Z"`
}

// ListCredentialsSuccessResponse represents the response for listing credentials
type ListCredentialsSuccessResponse struct {
	Credentials []CredentialsSuccessResponse `json:"credentials"`
}

// UploadCredentialResponse represents the response for uploading a credential file
type UploadCredentialResponse struct {
	DocumentID string `json:"documentId" example:"123e4567-e89b-12d3-a456-426614174000"`
	Status     string `json:"status" example:"Uploaded. Awaiting verification."`
}

type DeleteCredentialResponse struct {
	Message string `json:"message" example:"Operation successfully."`
}
