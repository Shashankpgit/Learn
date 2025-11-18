package service

import (
	"context"

	"app/src/constants"
	"app/src/model"
	"app/src/repository"
	"app/src/utils"
	"app/src/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ActorService defines the interface for actor business logic operations
type ActorService interface {
	RegisterActor(c *fiber.Ctx, req *validation.ApiRequest_RegistrationRequest) (*model.Actor, error)
	Login(c *fiber.Ctx, req *validation.ApiRequest_LoginRequest) (*AuthTokenResponse, error)
	Signout(c *fiber.Ctx, authUserID string) error
	UpdateActor(c *fiber.Ctx, actorID uuid.UUID, req *validation.ApiRequest_UpdateActorRequest) (*model.Actor, error)
	GetProfile(c *fiber.Ctx, actorID uuid.UUID) (*model.Actor, error)
	ForgotPassword(c *fiber.Ctx, req *validation.ApiRequest_ForgotPasswordRequest) error
	ResolveUniversalIdentifier(c *fiber.Ctx, req *validation.ApiRequest_ResolveRequest) (*model.Actor, error)
	GetUniversalIdentifier(actorID uuid.UUID) (string, error)
}

// actorService implements ActorService with constructor-based dependency injection
type actorService struct {
	log                  *logrus.Logger
	db                   *gorm.DB
	validate             *validator.Validate
	authService          AuthService
	actorRepo            repository.ActorRepository
	identifierRepo       repository.IdentifierRepository
	actorIntegrationRepo repository.ActorIntegrationRepository
}

// NewActorService creates a new actor service instance
func NewActorService(
	log *logrus.Logger,
	db *gorm.DB,
	validate *validator.Validate,
	authService AuthService,
	actorRepo repository.ActorRepository,
	identifierRepo repository.IdentifierRepository,
	actorIntegrationRepo repository.ActorIntegrationRepository,
) ActorService {
	return &actorService{
		log:                  log,
		db:                   db,
		validate:             validate,
		authService:          authService,
		actorRepo:            actorRepo,
		identifierRepo:       identifierRepo,
		actorIntegrationRepo: actorIntegrationRepo,
	}
}

// checkUniqueness is a generic helper for uniqueness validations
func (s *actorService) checkUniqueness(ctx context.Context, tx *gorm.DB, checkFn func(context.Context, *gorm.DB) (bool, error), errorMsg string, logMsg string) error {
	exists, err := checkFn(ctx, tx)
	if err != nil {
		s.log.Errorf("%s: %+v", logMsg, err)
		return err
	}
	if exists {
		return fiber.NewError(fiber.StatusConflict, errorMsg)
	}
	return nil
}

func (s *actorService) RegisterActor(c *fiber.Ctx, req *validation.ApiRequest_RegistrationRequest) (*model.Actor, error) {
	if err := s.validate.Struct(req.Request); err != nil {
		return nil, err
	}

	if err := req.Request.ValidateEntityTypeFields(); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var actor *model.Actor
	err := s.db.Transaction(func(tx *gorm.DB) error {
		ctx := c.Context()

		// Validate uniqueness constraints
		if err := s.validateUniquenessConstraints(ctx, tx, req); err != nil {
			return err
		}

		// Create actor
		actor = s.buildActorFromRequest(req)
		if err := s.actorRepo.Create(ctx, tx, actor); err != nil {
			s.log.Errorf("Failed to create actor: %+v", err)
			return err
		}

		// Create identifier if provided
		if req.Request.UniversalIdentifier != "" {
			if err := s.createIdentifier(ctx, tx, req.Request.UniversalIdentifier, actor.ActorID); err != nil {
				return err
			}
		}

		// Create auth user and integration
		if req.Request.UniversalIdentifier != "" {
			authUserID, err := s.authService.CreateUser(actor, req.Request.UniversalIdentifier, req.Request.Password)
			if err != nil {
				s.log.Errorf("Failed to create auth user: %+v", err)
				return err
			}
			s.log.Infof("Created auth user with ID: %s", authUserID)

			if authUserID != "" {
				if err := s.createActorIntegration(ctx, tx, actor.ActorID, authUserID); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return actor, err
}

// validateUniquenessConstraints validates all uniqueness constraints for actor creation
func (s *actorService) validateUniquenessConstraints(ctx context.Context, tx *gorm.DB, req *validation.ApiRequest_RegistrationRequest) error {
	// Check universal identifier
	if req.Request.UniversalIdentifier != "" {
		err := s.checkUniqueness(ctx, tx,
			func(ctx context.Context, tx *gorm.DB) (bool, error) {
				return s.identifierRepo.ExistsWithIdentifier(ctx, tx, req.Request.UniversalIdentifier)
			},
			constants.ErrUniversalIdentifierAlreadyInUse,
			"Failed to check universal identifier")
		if err != nil {
			return err
		}
	}

	// Check email
	if req.Request.Email != "" {
		err := s.checkUniqueness(ctx, tx,
			func(ctx context.Context, tx *gorm.DB) (bool, error) {
				return s.actorRepo.ExistsWithEmail(ctx, tx, req.Request.Email)
			},
			constants.ErrEmailAlreadyInUse,
			"Failed to check email")
		if err != nil {
			return err
		}
	}

	// Check master public key
	return s.checkUniqueness(ctx, tx,
		func(ctx context.Context, tx *gorm.DB) (bool, error) {
			return s.actorRepo.ExistsWithMasterPublicKey(ctx, tx, req.Request.MasterPublicKey)
		},
		constants.ErrMasterPublicKeyAlreadyInUse,
		"Failed to check master public key")
}

// buildActorFromRequest creates an actor model from the registration request
func (s *actorService) buildActorFromRequest(req *validation.ApiRequest_RegistrationRequest) *model.Actor {
	actor := &model.Actor{
		Email:             req.Request.Email,
		FirstName:         req.Request.FirstName,
		LastName:          req.Request.LastName,
		MasterPublicKey:   req.Request.MasterPublicKey,
		EntityType:        req.Request.EntityType,
		VerificationLevel: constants.VerificationLevelUnverified,
	}

	// Set optional fields
	actor.PhoneNumber = utils.StringPtr(req.Request.PhoneNumber)
	actor.Nationality = utils.StringPtr(req.Request.Nationality)
	actor.CountryOfResidence = utils.StringPtr(req.Request.CountryOfResidence)
	actor.CountryOfIncorporation = utils.StringPtr(req.Request.CountryOfIncorporation)

	return actor
}

// createIdentifier creates an identifier record
func (s *actorService) createIdentifier(ctx context.Context, tx *gorm.DB, identifier string, actorID uuid.UUID) error {
	id := &model.Identifier{
		Identifier: identifier,
		EntityType: constants.EntityTypeActor,
		EntityID:   actorID,
	}
	return s.identifierRepo.Create(ctx, tx, id)
}

// createActorIntegration creates an actor integration entry
func (s *actorService) createActorIntegration(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, authUserID string) error {
	integration := &model.ActorIntegration{
		ActorID:        actorID,
		Provider:       constants.KeycloakProviderName,
		ExternalUserID: authUserID,
	}
	return s.actorIntegrationRepo.Create(ctx, tx, integration)
}

func (s *actorService) Login(c *fiber.Ctx, req *validation.ApiRequest_LoginRequest) (*AuthTokenResponse, error) {
	if err := s.validate.Struct(req.Request); err != nil {
		return nil, err
	}

	s.log.Infof("Login attempt for username: %s", req.Request.Username)

	// Authenticate with auth provider
	tokenResp, err := s.authService.Login(req.Request.Username, req.Request.Password)
	if err != nil {
		s.log.Errorf("Authentication failed: %+v", err)
		return nil, err
	}

	// Find and validate actor
	actor, err := s.findActorByUsername(c.Context(), req.Request.Username)
	if err != nil {
		return nil, err
	}

	if actor.Email == "" {
		s.log.Errorf("Actor missing required email field: %s", actor.ActorID)
		return nil, fiber.NewError(fiber.StatusUnauthorized, constants.ErrInvalidActorAccount)
	}

	s.log.Infof("Login successful for actor: %s", actor.Email)
	return tokenResp, nil
}

// findActorByUsername finds an actor by username (universal identifier or email)
func (s *actorService) findActorByUsername(ctx context.Context, username string) (*model.Actor, error) {
	// Try by universal identifier first
	if identifier, err := s.identifierRepo.FindByValue(ctx, s.db, username); err == nil && identifier != nil {
		if actor, err := s.actorRepo.FindByID(ctx, s.db, identifier.EntityID); err == nil {
			return actor, nil
		}
	}

	// Try by email
	actor, err := s.actorRepo.FindByEmail(ctx, s.db, username)
	if err != nil {
		if utils.IsNotFoundError(err) {
			s.log.Errorf("Actor not found for username: %s", username)
			return nil, fiber.NewError(fiber.StatusUnauthorized, constants.ErrInvalidCredentials)
		}
		s.log.Errorf("Failed to find actor: %+v", err)
		return nil, err
	}

	s.log.Infof("Found actor for login: %s (ID: %s)", actor.Email, actor.ActorID)
	return actor, nil
}

func (s *actorService) Signout(c *fiber.Ctx, authUserID string) error {
	s.log.Infof("Initiating signout for auth user: %s", authUserID)

	if err := s.authService.Logout(authUserID); err != nil {
		s.log.Errorf("Failed to logout user: %+v", err)
		return err
	}

	s.log.Infof("Successfully signed out user: %s", authUserID)
	return nil
}

func (s *actorService) UpdateActor(c *fiber.Ctx, actorID uuid.UUID, req *validation.ApiRequest_UpdateActorRequest) (*model.Actor, error) {
	if err := s.validate.Struct(req.Request); err != nil {
		return nil, err
	}

	var actor *model.Actor
	err := s.db.Transaction(func(tx *gorm.DB) error {
		ctx := c.Context()

		var err error
		actor, err = s.actorRepo.FindByIDForUpdate(ctx, tx, actorID)
		if err != nil {
			return err
		}

		// Validate phone number uniqueness if provided
		if req.Request.PhoneNumber != nil && *req.Request.PhoneNumber != "" {
			err := s.checkUniqueness(ctx, tx,
				func(ctx context.Context, tx *gorm.DB) (bool, error) {
					return s.actorRepo.ExistsWithPhoneNumber(ctx, tx, *req.Request.PhoneNumber, actorID)
				},
				constants.ErrPhoneNumberAlreadyInUse,
				"Failed to check phone number")
			if err != nil {
				return err
			}
		}

		// Update fields
		s.updateActorFields(actor, req)

		if err := s.actorRepo.Update(ctx, tx, actor); err != nil {
			s.log.Errorf("Failed to update actor: %+v", err)
			return err
		}

		return nil
	})

	return actor, err
}

// updateActorFields updates actor fields from the request
func (s *actorService) updateActorFields(actor *model.Actor, req *validation.ApiRequest_UpdateActorRequest) {
	if req.Request.FirstName != "" {
		actor.FirstName = req.Request.FirstName
	}
	if req.Request.LastName != "" {
		actor.LastName = req.Request.LastName
	}
	if req.Request.PhoneNumber != nil {
		if *req.Request.PhoneNumber == "" {
			actor.PhoneNumber = nil
		} else {
			actor.PhoneNumber = req.Request.PhoneNumber
		}
	}
}

func (s *actorService) GetProfile(c *fiber.Ctx, actorID uuid.UUID) (*model.Actor, error) {
	return s.actorRepo.FindByID(c.Context(), s.db, actorID)
}

func (s *actorService) ForgotPassword(c *fiber.Ctx, req *validation.ApiRequest_ForgotPasswordRequest) error {
	if err := s.validate.Struct(req.Request); err != nil {
		return err
	}

	s.log.Infof("Password reset requested for email: %s", req.Request.Email)

	// Find actor by email
	actor, err := s.actorRepo.FindByEmail(c.Context(), s.db, req.Request.Email)
	if err != nil {
		if utils.IsNotFoundError(err) {
			s.log.Infof("Actor not found for email: %s", req.Request.Email)
			return fiber.NewError(fiber.StatusNotFound, constants.ErrResourceNotFound)
		}
		s.log.Errorf("Failed to find actor: %+v", err)
		return err
	}

	// Get auth integration
	integration, err := s.actorIntegrationRepo.FindByActorIDAndProvider(c.Context(), s.db, actor.ActorID, constants.KeycloakProviderName)
	if err != nil {
		if utils.IsNotFoundError(err) {
			s.log.Warnf("Auth integration not found for actor: %s", actor.ActorID)
			return fiber.NewError(fiber.StatusNotFound, constants.ErrResourceNotFound)
		}
		s.log.Errorf("Failed to find actor integration: %+v", err)
		return err
	}

	// Send password reset email
	if err := s.authService.ExecuteActionsEmail(integration.ExternalUserID, []string{"UPDATE_PASSWORD"}); err != nil {
		s.log.Errorf("Failed to send password reset email: %+v", err)
		return err
	}

	s.log.Infof("Password reset email sent successfully for actor: %s", actor.Email)
	return nil
}

func (s *actorService) ResolveUniversalIdentifier(c *fiber.Ctx, req *validation.ApiRequest_ResolveRequest) (*model.Actor, error) {
	if err := s.validate.Struct(req.Request); err != nil {
		return nil, err
	}

	ctx := c.Context()

	// Find identifier
	identifier, err := s.identifierRepo.FindByValue(ctx, s.db, req.Request.UniversalIdentifier)
	if err != nil {
		return nil, err
	}

	// Find actor
	actor, err := s.actorRepo.FindByID(ctx, s.db, identifier.EntityID)
	if err != nil {
		if utils.IsNotFoundError(err) {
			return nil, fiber.NewError(fiber.StatusNotFound, constants.ErrResourceNotFound)
		}
		return nil, err
	}

	return actor, nil
}

func (s *actorService) GetUniversalIdentifier(actorID uuid.UUID) (string, error) {
	identifier, err := s.identifierRepo.FindByActorID(context.Background(), s.db, actorID)
	if err != nil || identifier == nil {
		return "", err
	}
	return identifier.Identifier, nil
}
