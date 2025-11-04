package validation

import (
	"app/src/model"
)

// VerifiableCredential represents the verifiable credential structure for validation
type VerifiableCredential struct {
	Context           []string               `json:"@context" example:"https://www.w3.org/2018/credentials/v1"`
	ID                string                 `json:"id" example:"http://example.edu/credentials/3732"`
	Type              []string               `json:"type" example:"VerifiableCredential"`
	Issuer            string                 `json:"issuer" example:"did:example:issuer"`
	IssuanceDate      string                 `json:"issuanceDate" example:"2025-10-23T06:25:25.191Z"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	Proof             map[string]interface{} `json:"proof"`
}

// CredentialPayload represents the payload structure for validation
type CredentialPayload struct {
	VerifiableCredential VerifiableCredential `json:"verifiableCredential"`
	DocumentID           string               `json:"documentId" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// AddCredentialRequest represents the request structure for adding credentials
type AddCredentialRequest struct {
	VerificationType string            `json:"verificationType" example:"VC"`
	Payload          CredentialPayload `json:"payload"`
}

// ListCredentialsRequest represents the request for listing credentials
type ListCredentialsRequest struct{}

// GetCredentialRequest represents the request for getting a credential
type GetCredentialRequest struct {
	CredentialID string `json:"credentialId" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// DeleteCredentialRequest represents the request for deleting a credential
type DeleteCredentialRequest struct {
	CredentialID string `json:"credentialId" example:"123e4567-e89b-12d3-a456-426614174000"`
}

type AddCredentials struct {
	model.AddCredentialRequest
}

func ValidateAddCredentials(Req AddCredentialRequest) model.AddCredentialRequest {
	// Convert validation structures to model structures
	credentialSubject := make(map[string]interface{})
	for k, v := range Req.Payload.VerifiableCredential.CredentialSubject {
		credentialSubject[k] = v
	}

	proof := make(map[string]interface{})
	for k, v := range Req.Payload.VerifiableCredential.Proof {
		proof[k] = v
	}

	return model.AddCredentialRequest{
		VerificationType: Req.VerificationType,
		Payload: model.CredentialPayload{
			VerifiableCredential: model.VerifiableCredential{
				Context:           Req.Payload.VerifiableCredential.Context,
				ID:                Req.Payload.VerifiableCredential.ID,
				Type:              Req.Payload.VerifiableCredential.Type,
				Issuer:            Req.Payload.VerifiableCredential.Issuer,
				IssuanceDate:      Req.Payload.VerifiableCredential.IssuanceDate,
				CredentialSubject: credentialSubject,
				Proof:             proof,
			},
			DocumentID: Req.Payload.DocumentID,
		},
	}
}
