package utils

import (
	"app/src/constants"
	"app/src/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ResponseBuilder provides standardized methods for creating API responses
type ResponseBuilder struct{}

// NewResponseBuilder creates a new instance of ResponseBuilder
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{}
}

// extractRequestMetadata extracts id, ver, ts, msgid from the request body
func (rb *ResponseBuilder) extractRequestMetadata(c *fiber.Ctx) (id, ver, ts, msgid string) {
	// Try to extract from context locals first (if set by middleware)
	if reqID, ok := c.Locals("requestID").(string); ok && reqID != "" {
		id = reqID
	}
	if reqVer, ok := c.Locals("requestVersion").(string); ok && reqVer != "" {
		ver = reqVer
	}
	if reqTs, ok := c.Locals("requestTimestamp").(string); ok && reqTs != "" {
		ts = reqTs
	}
	if reqMsgID, ok := c.Locals("requestMsgID").(string); ok && reqMsgID != "" {
		msgid = reqMsgID
	}

	// If not in locals, use defaults
	if id == "" {
		id = constants.DefaultRequestID
	}
	if ver == "" {
		ver = constants.DefaultRequestVersion
	}
	if ts == "" {
		ts = time.Now().UTC().Format(time.RFC3339)
	}
	if msgid == "" {
		msgid = constants.DefaultMsgID
	}

	return id, ver, ts, msgid
}

// Success sends a successful response with the given data
func (rb *ResponseBuilder) Success(c *fiber.Ctx, statusCode int, data interface{}) error {
	id, ver, ts, msgid := rb.extractRequestMetadata(c)

	apiResponse := response.ApiResponse{
		ID:  id,
		Ver: ver,
		Ts:  ts,
		Params: response.ApiResponseParams{
			MsgID:  msgid,
			Status: constants.StatusSuccessful,
		},
		Response: data,
	}

	return c.Status(statusCode).JSON(apiResponse)
}

// SuccessWithMetadata sends a successful response with explicit metadata
func (rb *ResponseBuilder) SuccessWithMetadata(c *fiber.Ctx, statusCode int, reqID, reqVer, reqTs, msgid string, data interface{}) error {
	apiResponse := response.ApiResponse{
		ID:  reqID,
		Ver: reqVer,
		Ts:  reqTs,
		Params: response.ApiResponseParams{
			MsgID:  msgid,
			Status: constants.StatusSuccessful,
		},
		Response: data,
	}

	return c.Status(statusCode).JSON(apiResponse)
}

// Error sends an error response with the given status code and message
func (rb *ResponseBuilder) Error(c *fiber.Ctx, statusCode int, errorMessage string) error {
	id, ver, ts, msgid := rb.extractRequestMetadata(c)
	errorCode := mapStatusToErrorCode(statusCode)

	apiResponse := response.ApiResponse_Error{
		ApiResponse: response.ApiResponse{
			ID:  id,
			Ver: ver,
			Ts:  ts,
			Params: response.ApiResponseParams{
				MsgID:  msgid,
				Status: constants.StatusFailed,
				Err:    errorCode,
				ErrMsg: errorMessage,
			},
			Response: struct{}{},
		},
	}

	return c.Status(statusCode).JSON(apiResponse)
}

// ErrorWithMetadata sends an error response with explicit metadata
func (rb *ResponseBuilder) ErrorWithMetadata(c *fiber.Ctx, statusCode int, reqID, reqVer, reqTs, msgid string, errorMessage string) error {
	errorCode := mapStatusToErrorCode(statusCode)

	apiResponse := response.ApiResponse_Error{
		ApiResponse: response.ApiResponse{
			ID:  reqID,
			Ver: reqVer,
			Ts:  reqTs,
			Params: response.ApiResponseParams{
				MsgID:  msgid,
				Status: constants.StatusFailed,
				Err:    errorCode,
				ErrMsg: errorMessage,
			},
			Response: struct{}{},
		},
	}

	return c.Status(statusCode).JSON(apiResponse)
}

// ValidationError sends a validation error response with field-specific errors
func (rb *ResponseBuilder) ValidationError(c *fiber.Ctx, validationErrors map[string]string) error {
	id, ver, ts, msgid := rb.extractRequestMetadata(c)

	// Format validation errors as a readable string
	errMsg := constants.MsgBadRequest
	for field, msg := range validationErrors {
		errMsg += "; " + field + ": " + msg
	}

	apiResponse := response.ApiResponse_Error{
		ApiResponse: response.ApiResponse{
			ID:  id,
			Ver: ver,
			Ts:  ts,
			Params: response.ApiResponseParams{
				MsgID:  msgid,
				Status: constants.StatusFailed,
				Err:    constants.ErrCodeBadRequest,
				ErrMsg: errMsg,
			},
			Response: struct{}{},
		},
	}

	return c.Status(fiber.StatusBadRequest).JSON(apiResponse)
}

// Created sends a 201 Created response
func (rb *ResponseBuilder) Created(c *fiber.Ctx, data interface{}) error {
	return rb.Success(c, fiber.StatusCreated, data)
}

// CreatedWithMetadata sends a 201 Created response with explicit metadata
func (rb *ResponseBuilder) CreatedWithMetadata(c *fiber.Ctx, reqID, reqVer, reqTs, msgid string, data interface{}) error {
	return rb.SuccessWithMetadata(c, fiber.StatusCreated, reqID, reqVer, reqTs, msgid, data)
}

// Accepted sends a 202 Accepted response
func (rb *ResponseBuilder) Accepted(c *fiber.Ctx, data interface{}) error {
	return rb.Success(c, fiber.StatusAccepted, data)
}

// AcceptedWithMetadata sends a 202 Accepted response with explicit metadata
func (rb *ResponseBuilder) AcceptedWithMetadata(c *fiber.Ctx, reqID, reqVer, reqTs, msgid string, data interface{}) error {
	return rb.SuccessWithMetadata(c, fiber.StatusAccepted, reqID, reqVer, reqTs, msgid, data)
}

// OK sends a 200 OK response
func (rb *ResponseBuilder) OK(c *fiber.Ctx, data interface{}) error {
	return rb.Success(c, fiber.StatusOK, data)
}

// OKWithMetadata sends a 200 OK response with explicit metadata
func (rb *ResponseBuilder) OKWithMetadata(c *fiber.Ctx, reqID, reqVer, reqTs, msgid string, data interface{}) error {
	return rb.SuccessWithMetadata(c, fiber.StatusOK, reqID, reqVer, reqTs, msgid, data)
}

// NoContent sends a 204 No Content response
func (rb *ResponseBuilder) NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// BadRequest sends a 400 Bad Request response
func (rb *ResponseBuilder) BadRequest(c *fiber.Ctx, message string) error {
	return rb.Error(c, fiber.StatusBadRequest, message)
}

// Unauthorized sends a 401 Unauthorized response
func (rb *ResponseBuilder) Unauthorized(c *fiber.Ctx, message string) error {
	return rb.Error(c, fiber.StatusUnauthorized, message)
}

// Forbidden sends a 403 Forbidden response
func (rb *ResponseBuilder) Forbidden(c *fiber.Ctx, message string) error {
	return rb.Error(c, fiber.StatusForbidden, message)
}

// NotFound sends a 404 Not Found response
func (rb *ResponseBuilder) NotFound(c *fiber.Ctx, message string) error {
	return rb.Error(c, fiber.StatusNotFound, message)
}

// Conflict sends a 409 Conflict response
func (rb *ResponseBuilder) Conflict(c *fiber.Ctx, message string) error {
	return rb.Error(c, fiber.StatusConflict, message)
}

// InternalServerError sends a 500 Internal Server Error response
func (rb *ResponseBuilder) InternalServerError(c *fiber.Ctx, message string) error {
	return rb.Error(c, fiber.StatusInternalServerError, message)
}

// mapStatusToErrorCode maps HTTP status code to error code
func mapStatusToErrorCode(statusCode int) string {
	switch statusCode {
	case fiber.StatusBadRequest:
		return constants.ErrCodeBadRequest
	case fiber.StatusUnauthorized:
		return constants.ErrCodeUnauthorized
	case fiber.StatusNotFound:
		return constants.ErrCodeNotFound
	case fiber.StatusConflict:
		return constants.ErrCodeConflict
	case fiber.StatusInternalServerError:
		return constants.ErrCodeInternalServerError
	default:
		return constants.ErrCodeInternalServerError
	}
}

// Global instance for convenience
var Response = NewResponseBuilder()
