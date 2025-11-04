package service

import (
	"app/src/constants"
	"app/src/model"
	"app/src/repository"
	"app/src/utils"
	"app/src/validation"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CredentialsService defines the interface for credentials business logic operations
type CredentialsService interface {
	AddCredential(c *fiber.Ctx, req *validation.AddCredentials) (*model.Token, error)
	ListCredentials(c *fiber.Ctx) ([]model.Token, error)
	GetCredential(c *fiber.Ctx, credentialID string) (*model.Token, error)
	DeleteCredential(c *fiber.Ctx, credentialID string) error
}

// credentialsService implements CredentialsService with constructor-based dependency injection
type credentialsService struct {
	log             *logrus.Logger
	db              *gorm.DB
	validate        *validator.Validate
	credentialsRepo repository.CredentialsRepository
	documentRepo    repository.DocumentRepository
}

// NewCredentialsService creates a new credentials service instance
func NewCredentialsService(
	log *logrus.Logger,
	db *gorm.DB,
	validate *validator.Validate,
	credentialsRepo repository.CredentialsRepository,
	documentRepo repository.DocumentRepository,
) CredentialsService {
	return &credentialsService{
		log:             log,
		db:              db,
		validate:        validate,
		credentialsRepo: credentialsRepo,
		documentRepo:    documentRepo,
	}
}

func (s *credentialsService) AddCredential(c *fiber.Ctx, req *validation.AddCredentials) (*model.Token, error) {
	if err := s.validate.Struct(req.AddCredentialRequest); err != nil {
		return nil, err
	}

	actorID := c.Locals("actorID")
	if actorID == nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, constants.ErrUnauthorized)
	}

	actorUUID, ok := actorID.(uuid.UUID)
	if !ok {
		s.log.Error("actorID type assertion failed")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "invalid actor ID type")
	}

	// Validate and parse document ID
	if req.AddCredentialRequest.Payload.DocumentID == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "documentId is required")
	}

	documentID, err := utils.ParseUUID(req.AddCredentialRequest.Payload.DocumentID, "document")
	if err != nil {
		return nil, err
	}

	// Verify document exists
	if _, err = s.documentRepo.FindByID(c.Context(), s.db, documentID); err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, constants.ErrDocumentNotFound)
	}

	// Build token
	token := s.buildTokenFromRequest(req)
	token.AccountID = actorUUID

	// Save token - single operation doesn't need transaction
	if err := s.credentialsRepo.Create(c.Context(), s.db, token); err != nil {
		s.log.Errorf("%s: %+v", constants.ErrFailedToCreateToken, err)
		return nil, err
	}

	return token, nil
}

// buildTokenFromRequest creates a token from the credential request
func (s *credentialsService) buildTokenFromRequest(req *validation.AddCredentials) *model.Token {
	vc := req.AddCredentialRequest.Payload.VerifiableCredential

	metadata := map[string]interface{}{
		"verifiableCredential": map[string]interface{}{
			"@context":          vc.Context,
			"id":                vc.ID,
			"type":              vc.Type,
			"issuer":            vc.Issuer,
			"issuanceDate":      vc.IssuanceDate,
			"credentialSubject": vc.CredentialSubject,
			"proof":             vc.Proof,
		},
		"verificationType": req.AddCredentialRequest.VerificationType,
	}

	return &model.Token{
		TokenID:       uuid.New(),
		TokenType:     req.AddCredentialRequest.VerificationType,
		IssuerDID:     vc.Issuer,
		TokenStandard: constants.TokenStandardVC,
		Status:        constants.StatusPending,
		Metadata:      datatypes.JSONMap(metadata),
	}
}

func (s *credentialsService) ListCredentials(c *fiber.Ctx) ([]model.Token, error) {
	tokens, err := s.credentialsRepo.FindAll(c.Context(), s.db)
	if err != nil {
		s.log.Errorf("Failed to retrieve credentials: %+v", err)
		return nil, err
	}

	if len(tokens) == 0 {
		return []model.Token{}, nil
	}

	return tokens, nil
}

func (s *credentialsService) GetCredential(c *fiber.Ctx, credentialID string) (*model.Token, error) {
	tokenID, err := utils.ParseUUID(credentialID, "credential")
	if err != nil {
		return nil, err
	}

	token, err := s.credentialsRepo.FindByID(c.Context(), s.db, tokenID)
	if err != nil {
		s.log.Errorf("Failed to retrieve credential: %+v", err)
		return nil, err
	}

	return token, nil
}

func (s *credentialsService) DeleteCredential(c *fiber.Ctx, credentialID string) error {
	tokenID, err := utils.ParseUUID(credentialID, "credential")
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.credentialsRepo.Delete(c.Context(), tx, tokenID); err != nil {
			s.log.Errorf("Failed to delete credential: %+v", err)
			return err
		}
		return nil
	})
}
