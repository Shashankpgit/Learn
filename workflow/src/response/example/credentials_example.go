package example

// ParamsExample models the `params` envelope shown in the user's sample JSON.
// It includes message id, status and error fields commonly used in error replies.
type ParamsExample struct {
	MsgID  string `json:"msgid" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Status string `json:"status" example:"failed"`
	Err    string `json:"err" example:"IDENTIFIER_TAKEN"`
	ErrMsg string `json:"errmsg" example:"The universal identifier 'alice@finternet' is already taken."`
}

// Per-error params examples so Swagger can show different `err`/`errmsg`
// values for each response type instead of the same generic params.
type ParamsBadRequestExample struct {
	MsgID  string `json:"msgid" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Status string `json:"status" example:"failed"`
	Err    string `json:"err" example:"INVALID_REQUEST"`
	ErrMsg string `json:"errmsg" example:"Invalid request body"`
}

type NotFoundExample struct {
	MsgID  string `json:"msgid" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Status string `json:"status" example:"failed"`
	Err    string `json:"err" example:"CREDENTIAL_NOT_FOUND"`
	ErrMsg string `json:"errmsg" example:"Credential not found"`
}

type DocumentNotFoundExample struct {
	MsgID  string `json:"msgid" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Status string `json:"status" example:"failed"`
	Err    string `json:"err" example:"DOCUMENT_NOT_FOUND"`
	ErrMsg string `json:"errmsg" example:"Document ID not found"`
}

type ParamsInternalExample struct {
	MsgID  string `json:"msgid" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Status string `json:"status" example:"failed"`
	Err    string `json:"err" example:"INTERNAL_SERVER_ERROR"`
	ErrMsg string `json:"errmsg" example:"Internal server error"`
}

type ParamsConflictExample struct {
	MsgID  string `json:"msgid" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Status string `json:"status" example:"failed"`
	Err    string `json:"err" example:"IDENTIFIER_TAKEN"`
	ErrMsg string `json:"errmsg" example:"The universal identifier 'alice@finternet' is already taken."`
}

type ParamsUnprocessableEntityExample struct {
	MsgID  string `json:"msgid" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Status string `json:"status" example:"failed"`
	Err    string `json:"err" example:"UNPROCESSABLE_ENTITY"`
	ErrMsg string `json:"errmsg" example:"Unable to process the contained instructions"`
}

type BadRequestExample struct {
	MsgID  string `json:"msgid" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Status string `json:"status" example:"failed"`
	Err    string `json:"err" example:"BAD_REQUEST"`
	ErrMsg string `json:"errmsg" example:"Invalid request body"`
}

type UnauthorizedExample struct {
	MsgID  string `json:"msgid" example:"a1b2c3d4-e5f6-7890-1234-567890abcdef"`
	Status string `json:"status" example:"failed"`
	Err    string `json:"err" example:"UNAUTHORIZED"`
	ErrMsg string `json:"errmsg" example:"Unauthorized. The JWT is missing, invalid, or expired"`
}

// A sample envelope matching the exact JSON structure provided by the user.
// Useful for swagger examples or documentation.
type ErrorEnvelope[T any] struct {
	ID       string      `json:"id" example:"api_id"`
	Ver      string      `json:"ver" example:"v1"`
	Ts       string      `json:"ts" example:"2025-09-08T12:00:00Z"`
	Params   T           `json:"params"`
	Response interface{} `json:"response" swaggertype:"object"`
}
