package controller

import (
	"app/src/constants"
	"app/src/middleware"
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/utils"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ActorController struct {
	ActorService    service.ActorService
	ResponseBuilder *utils.ResponseBuilder
}

func NewActorController(
	actorService service.ActorService,
	responseBuilder *utils.ResponseBuilder,
) *ActorController {
	return &ActorController{
		ActorService:    actorService,
		ResponseBuilder: responseBuilder,
	}
}

// getActorID extracts the authenticated actor ID from context
func (a *ActorController) getActorID(c *fiber.Ctx) (uuid.UUID, error) {
	actorID, ok := c.Locals("actorID").(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, constants.ErrUnauthorized)
	}
	return actorID, nil
}

// getAuthUserID extracts the auth user ID from JWT claims
func (a *ActorController) getAuthUserID(c *fiber.Ctx) (string, error) {
	claims, ok := c.Locals("auth_claims").(*middleware.AuthClaims)
	if !ok || claims == nil || claims.Sub == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, constants.ErrUnauthorized)
	}
	return claims.Sub, nil
}

// getUniversalIdentifier retrieves the universal identifier for an actor
func (a *ActorController) getUniversalIdentifier(actorID uuid.UUID) string {
	identifier, _ := a.ActorService.GetUniversalIdentifier(actorID)
	return identifier
}

// @Tags         Actor
// @Summary      Register a new actor
// @Description  Creates a new user actor, generates and stores the reproducible DID from the master public key, and links them.
// @Accept       json
// @Produce      json
// @Param        request  body  validation.ApiRequest_RegistrationRequest  true  "Request body"
// @Router       /v1/actor/create [post]
// @Success      201  {object}  response.ApiResponse_RegistrationSuccess
// @Failure      400  {object}  response.ApiResponse_Error  "Bad request, such as a missing required field"
// @Failure      409  {object}  response.ApiResponse_Error  "Conflict, an identifier is already in use"
func (a *ActorController) RegisterActor(c *fiber.Ctx) error {
	var req validation.ApiRequest_RegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	actor, err := a.ActorService.RegisterActor(c, &req)
	if err != nil {
		return err
	}

	responseData := response.RegistrationSuccessResponse{
		DID:     actor.DID,
		Message: constants.MsgActorRegisteredSuccessfully,
	}

	return a.ResponseBuilder.CreatedWithMetadata(c, req.ID, req.Ver, req.Ts, *req.Params.MsgID, responseData)
}

// @Tags         Actor
// @Summary      Login actor
// @Description  Authenticates an actor and returns JWT tokens
// @Accept       json
// @Produce      json
// @Param        request  body  validation.ApiRequest_LoginRequest  true  "Request body"
// @Router       /v1/actor/login [post]
// @Success      200  {object}  response.ApiResponse  "Login successful"
// @Failure      400  {object}  response.ApiResponse_Error  "Bad request"
// @Failure      401  {object}  response.ApiResponse_Error  "Invalid credentials"
func (a *ActorController) Login(c *fiber.Ctx) error {
	var req validation.ApiRequest_LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	authToken, err := a.ActorService.Login(c, &req)
	if err != nil {
		return err
	}

	responseData := response.AuthResponse{
		AccessToken: authToken.AccessToken,
		TokenType:   authToken.TokenType,
		ExpiresIn:   authToken.ExpiresIn,
	}

	return a.ResponseBuilder.OKWithMetadata(c, req.ID, req.Ver, req.Ts, *req.Params.MsgID, responseData)
}

// @Tags         Actor
// @Summary      Update actor profile
// @Description  Allows an authenticated user to update their profile information, such as name and phone number. The actor to update is determined by the session token.
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  validation.ApiRequest_UpdateActorRequest  true  "Request body"
// @Router       /v1/actor/update [post]
// @Success      200  {object}  response.ApiResponse_ActorProfile
// @Failure      400  {object}  response.ApiResponse_Error  "Bad request"
// @Failure      401  {object}  response.ApiResponse_Error  "Unauthorized"
// @Failure      409  {object}  response.ApiResponse_Error  "Conflict, the new phone number is already in use"
func (a *ActorController) UpdateActor(c *fiber.Ctx) error {
	var req validation.ApiRequest_UpdateActorRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	actorID, err := a.getActorID(c)
	if err != nil {
		return err
	}

	actor, err := a.ActorService.UpdateActor(c, actorID, &req)
	if err != nil {
		return err
	}

	responseData := a.buildActorProfileResponse(actor)
	return a.ResponseBuilder.OKWithMetadata(c, req.ID, req.Ver, req.Ts, *req.Params.MsgID, responseData)
}

// @Tags         Actor
// @Summary      Get current actor's profile
// @Description  Retrieves the full actor profile for the authenticated actor.
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  validation.ApiRequest_Empty  true  "Request body"
// @Router       /v1/actor/getProfile [post]
// @Success      200  {object}  response.ApiResponse_ActorProfile
// @Failure      401  {object}  response.ApiResponse_Error  "Unauthorized"
func (a *ActorController) GetProfile(c *fiber.Ctx) error {
	var req validation.ApiRequest_Empty
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	actorID, err := a.getActorID(c)
	if err != nil {
		return err
	}

	actor, err := a.ActorService.GetProfile(c, actorID)
	if err != nil {
		return err
	}

	responseData := a.buildActorProfileResponse(actor)
	return a.ResponseBuilder.OKWithMetadata(c, req.ID, req.Ver, req.Ts, *req.Params.MsgID, responseData)
}

// @Tags         Actor
// @Summary      Actor signout
// @Description  Logs out the user from auth provider and revokes all their active sessions.
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  validation.ApiRequest_Empty  true  "Request body"
// @Router       /v1/actor/signout [post]
// @Success      200  {object}  response.ApiResponse_Success
// @Failure      401  {object}  response.ApiResponse_Error  "Unauthorized"
func (a *ActorController) Signout(c *fiber.Ctx) error {
	var req validation.ApiRequest_Empty
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	authUserID, err := a.getAuthUserID(c)
	if err != nil {
		return err
	}

	if err := a.ActorService.Signout(c, authUserID); err != nil {
		return err
	}

	responseData := response.SuccessResponse{
		Message: constants.MsgSignoutSuccessful,
	}

	return a.ResponseBuilder.OKWithMetadata(c, req.ID, req.Ver, req.Ts, *req.Params.MsgID, responseData)
}

// @Tags         Actor
// @Summary      Forgot password
// @Description  Initiates the password reset flow for a given email address. A reset link will be sent to the user's email.
// @Accept       json
// @Produce      json
// @Param        request  body  validation.ApiRequest_ForgotPasswordRequest  true  "Request body"
// @Router       /v1/actor/forgotPassword [post]
// @Success      200  {object}  response.ApiResponse_Success
// @Failure      404  {object}  response.ApiResponse_Error  "Not found"
func (a *ActorController) ForgotPassword(c *fiber.Ctx) error {
	var req validation.ApiRequest_ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	if err := a.ActorService.ForgotPassword(c, &req); err != nil {
		return err
	}

	responseData := response.SuccessResponse{
		Message: constants.MsgPasswordResetInstructions,
	}

	return a.ResponseBuilder.OKWithMetadata(c, req.ID, req.Ver, req.Ts, *req.Params.MsgID, responseData)
}

// @Tags         Actor
// @Summary      Resolve a universal identifier
// @Description  Publicly resolves a universal identifier to its master public key and DID.
// @Accept       json
// @Produce      json
// @Param        request  body  validation.ApiRequest_ResolveRequest  true  "Request body"
// @Router       /v1/actor/resolve [post]
// @Success      200  {object}  response.ApiResponse_ResolveResponse
// @Failure      404  {object}  response.ApiResponse_Error  "Not found"
func (a *ActorController) ResolveUniversalIdentifier(c *fiber.Ctx) error {
	var req validation.ApiRequest_ResolveRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	actor, err := a.ActorService.ResolveUniversalIdentifier(c, &req)
	if err != nil {
		return err
	}

	responseData := response.ResolveResponse{
		UniversalIdentifier: req.Request.UniversalIdentifier,
		MasterPublicKey:     actor.MasterPublicKey,
		DID:                 actor.DID,
	}

	return a.ResponseBuilder.OKWithMetadata(c, req.ID, req.Ver, req.Ts, *req.Params.MsgID, responseData)
}

// buildActorProfileResponse builds an actor profile response
func (a *ActorController) buildActorProfileResponse(actor *model.Actor) response.ActorProfile {
	return response.ActorProfile{
		DID:                    actor.DID,
		UniversalIdentifier:    a.getUniversalIdentifier(actor.ActorID),
		Email:                  actor.Email,
		FirstName:              actor.FirstName,
		LastName:               actor.LastName,
		PhoneNumber:            actor.PhoneNumber,
		MasterPublicKey:        actor.MasterPublicKey,
		VerificationLevel:      actor.VerificationLevel,
		EntityType:             actor.EntityType,
		Nationality:            actor.Nationality,
		CountryOfResidence:     actor.CountryOfResidence,
		CountryOfIncorporation: actor.CountryOfIncorporation,
	}
}
