package controller

import (
	"app/src/adapter"
	"app/src/constants"
	"app/src/model"
	"app/src/repository"
	"app/src/response"
	_ "app/src/response/example"
	"app/src/service"
	"app/src/utils"
	"app/src/validation"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CredentialController handles credential-related HTTP requests
type CredentialController struct {
	credentialsService service.CredentialsService
	storageFactory     *adapter.StorageFactory
	storageProvider    adapter.StorageProvider
	storageOnce        sync.Once
	storageErr         error
	db                 *gorm.DB
	documentRepo       repository.DocumentRepository
	responseBuilder    *utils.ResponseBuilder
}

// NewCredentialsController creates a new credentials controller
// Storage provider is initialized lazily on first use to avoid startup failures
func NewCredentialsController(
	credentialsService service.CredentialsService,
	storageFactory *adapter.StorageFactory,
	db *gorm.DB,
	documentRepo repository.DocumentRepository,
	responseBuilder *utils.ResponseBuilder,
) *CredentialController {
	return &CredentialController{
		credentialsService: credentialsService,
		storageFactory:     storageFactory,
		db:                 db,
		documentRepo:       documentRepo,
		responseBuilder:    responseBuilder,
	}
}

// getStorageProvider returns the storage provider, initializing it lazily on first call
func (cc *CredentialController) getStorageProvider() (adapter.StorageProvider, error) {
	cc.storageOnce.Do(func() {
		cc.storageProvider, cc.storageErr = cc.storageFactory.NewProvider()
		if cc.storageErr != nil {
			cc.storageErr = fmt.Errorf("%s: %w", constants.ErrFailedToCreateStorageProvider, cc.storageErr)
		}
	})
	return cc.storageProvider, cc.storageErr
}

// @Tags         Credentials
// @Summary      Add a credential for verification
// @Description  Submits a credential (VC or manual document ref) for asynchronous verification, creating a credential token in 'Pending' state.
// @Produce      json
// @Param        request body  response.Request[validation.AddCredentialRequest]  true  "Request body"
// @Router       /credentials/add [post]
// @Success      202  {object}  response.Response[response.AddCredentialSuccessResponse]  "Credential accepted for verification"
// @Failure      400  {object}  example.ErrorEnvelope[example.BadRequestExample]  "Invalid request"
// @Failure      401  {object}  example.ErrorEnvelope[example.UnauthorizedExample]  "Unauthorized. JWT missing, invalid, or expired."
// @Failure      404  {object}  example.ErrorEnvelope[example.DocumentNotFoundExample]  "Document ID not found"
func (cc *CredentialController) AddCredential(c *fiber.Ctx) error {
	var req response.Request[validation.AddCredentialRequest]
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	credentials := &validation.AddCredentials{
		AddCredentialRequest: validation.ValidateAddCredentials(req.Request),
	}

	token, err := cc.credentialsService.AddCredential(c, credentials)
	if err != nil {
		return err
	}

	payload := response.AddCredentialSuccessResponse{
		CredentialID: token.TokenID.String(),
		Status:       constants.StatusPending,
		Message:      constants.MsgCredentialSubmittedSuccessfully,
	}

	return cc.responseBuilder.AcceptedWithMetadata(c, req.ID, req.Ver, req.Ts, *req.Params.MsgID, payload)
}

// @Tags         Credentials
// @Summary      List credentials
// @Description  List all credentials for a user
// @Produce      json
// @Param        request body  response.Request[validation.ListCredentialsRequest]  true  "Request body"
// @Router       /credentials/list [post]
// @Success      200 {object} response.Response[response.ListCredentialsSuccessResponse]
// @Failure      400  {object}  example.ErrorEnvelope[example.BadRequestExample]  "Invalid request body"
// @Failure      401  {object}  example.ErrorEnvelope[example.UnauthorizedExample]  "Unauthorized. JWT missing, invalid, or expired."
// @Failure      500  {object}  example.ErrorEnvelope[example.ParamsInternalExample]  "Internal server error"
func (cc *CredentialController) ListCredentials(c *fiber.Ctx) error {
	var req response.Request[validation.ListCredentialsRequest]
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	tokens, err := cc.credentialsService.ListCredentials(c)
	if err != nil {
		return err
	}

	payload := make([]response.CredentialsSuccessResponse, 0, len(tokens))
	for _, token := range tokens {
		payload = append(payload, response.CredentialsSuccessResponse{
			CredentialID: token.TokenID.String(),
			Type:         constants.CredentialTypeVC,
			Status:       constants.StatusPending,
			SubmittedAt:  token.CreatedAt.Format(time.RFC3339),
		})
	}

	return cc.responseBuilder.OKWithMetadata(c, req.ID, req.Ver, req.Ts, *req.Params.MsgID, payload)
}

// @Tags         Credentials
// @Summary      Get credential
// @Description  Get credential for a user
// @Produce      json
// @Param        request body  response.Request[validation.GetCredentialRequest]  true  "Request body"
// @Router       /credentials/get [post]
// @Success      200  {object}  response.Response[response.CredentialsSuccessResponse]
// @Failure      400  {object}  example.ErrorEnvelope[example.BadRequestExample]  "Invalid request body"
// @Failure      401  {object}  example.ErrorEnvelope[example.UnauthorizedExample]  "Unauthorized. JWT missing, invalid, or expired."
// @Failure      404  {object}  example.ErrorEnvelope[example.NotFoundExample]  "Credential not found"
// @Failure      500  {object}  example.ErrorEnvelope[example.ParamsInternalExample]  "Internal server error"
func (cc *CredentialController) GetCredential(c *fiber.Ctx) error {
	var req response.Request[validation.GetCredentialRequest]
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	token, err := cc.credentialsService.GetCredential(c, req.Request.CredentialID)
	if err != nil {
		return err
	}

	payload := response.CredentialsSuccessResponse{
		CredentialID: req.Request.CredentialID,
		Type:         constants.CredentialTypeVC,
		Status:       constants.StatusPending,
		SubmittedAt:  token.CreatedAt.Format(time.RFC3339),
	}

	return cc.responseBuilder.OKWithMetadata(c, req.ID, req.Ver, req.Ts, *req.Params.MsgID, payload)
}

// @Tags         Credentials
// @Summary      Delete credential
// @Description  Delete credential for a user
// @Produce      json
// @Param        request body  response.Request[validation.DeleteCredentialRequest]  true  "Request body"
// @Router       /credentials/delete [post]
// @Success      200  {object}  response.Response[response.DeleteCredentialResponse]
// @Failure      400  {object}  example.ErrorEnvelope[example.BadRequestExample]  "Invalid request body"
// @Failure      401  {object}  example.ErrorEnvelope[example.UnauthorizedExample]  "Unauthorized. JWT missing, invalid, or expired."
// @Failure      404  {object}  example.ErrorEnvelope[example.NotFoundExample]  "Credential not found"
// @Failure      500  {object}  example.ErrorEnvelope[example.ParamsInternalExample]  "Internal server error"
func (cc *CredentialController) DeleteCredential(c *fiber.Ctx) error {
	var req response.Request[validation.DeleteCredentialRequest]
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrInvalidRequestBody)
	}

	if err := cc.credentialsService.DeleteCredential(c, req.Request.CredentialID); err != nil {
		return err
	}

	payload := map[string]string{
		"message": constants.MsgOperationSuccessful,
	}

	return cc.responseBuilder.OKWithMetadata(c, req.ID, req.Ver, req.Ts, *req.Params.MsgID, payload)
}

// @Tags         Credentials
// @Summary      Upload a document for verification
// @Description  Uploads a document (e.g., passport) for manual verification. Returns a document ID used in /credentials/add.
// @Accept       multipart/form-data
// @Produce      json
// @Param        document formData file true "Document to upload"
// @Router       /credentials/upload [post]
// @Success      201  {object}  response.Response[response.UploadCredentialResponse]  "Document uploaded successfully"
// @Failure      400  {object}  example.ErrorEnvelope[example.BadRequestExample]  "Invalid request"
// @Failure      401  {object}  example.ErrorEnvelope[example.UnauthorizedExample]  "Unauthorized. JWT missing, invalid, or expired."
func (cc *CredentialController) UploadFile(c *fiber.Ctx) error {
	// Get storage provider (lazy initialization)
	storageProvider, err := cc.getStorageProvider()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, 
			fmt.Sprintf("Storage provider unavailable: %v", err))
	}

	actorID := uuid.Must(uuid.NewV7())

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, constants.ErrFileRequired)
	}

	fileReader, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, constants.ErrFailedToOpenFile)
	}
	defer fileReader.Close()

	// Build storage key and upload file
	storageKey, fileExt := cc.buildStorageKey(file.Filename)
	opts := &adapter.UploadOptions{
		ContentType: file.Header.Get(constants.HTTPHeaderContentType),
		Metadata: map[string]string{
			"original-filename": file.Filename,
			"uploaded-at":       time.Now().Format(time.RFC3339),
			"actor-id":          actorID.String(),
		},
	}

	if _, err = storageProvider.Upload(c.Context(), storageKey, fileReader, file.Size, opts); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, 
			fmt.Sprintf("%s: %v", constants.ErrFailedToUploadFile, err))
	}

	// Create document record
	document := &model.Document{
		DocumentID:  uuid.Must(uuid.NewV7()),
		AccountID:   actorID,
		FileName:    file.Filename,
		StoragePath: storageKey,
		MimeType:    &fileExt,
		UploadedAt:  time.Now(),
	}

	if err := cc.db.Transaction(func(tx *gorm.DB) error {
		return cc.documentRepo.Create(c.Context(), tx, document)
	}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, 
			fmt.Sprintf("Failed to save document record: %v", err))
	}

	payload := response.UploadCredentialResponse{
		DocumentID: document.DocumentID.String(),
		Status:     constants.MsgUploadedAwaitingVerification,
	}

	return cc.responseBuilder.CreatedWithMetadata(c, 
		constants.DefaultRequestID, 
		constants.DefaultRequestVersion, 
		time.Now().UTC().Format(time.RFC3339), 
		constants.DefaultMsgID, 
		payload)
}

// buildStorageKey generates a unique storage key from filename
func (cc *CredentialController) buildStorageKey(filename string) (storageKey, fileExt string) {
	fileExt = strings.TrimPrefix(filepath.Ext(filename), ".")
	if fileExt == "" {
		fileExt = "bin"
	}

	baseName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	storageKey = fmt.Sprintf("/%s_%d.%s", baseName, time.Now().UnixNano(), fileExt)
	return
}
