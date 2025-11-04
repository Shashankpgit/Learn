package constants

// Error Messages and Codes
const (
	// Validation Error Messages
	ErrInvalidRequestBody                        = "Invalid request body"
	ErrInvalidCredentials                        = "Invalid credentials"
	ErrInvalidActorAccount                       = "Invalid actor account"
	ErrUnauthorized                              = "Unauthorized. The JWT is missing, invalid, or expired"
	ErrResourceNotFound                          = "The requested resource was not found."
	ErrActorNotFound                             = "Actor not found"
	ErrEmailAlreadyInUse                         = "Email is already in use"
	ErrMasterPublicKeyAlreadyInUse               = "Master public key is already in use"
	ErrPhoneNumberAlreadyInUse                   = "Phone number is already in use"
	ErrUniversalIdentifierAlreadyInUse           = "Universal identifier is already in use"
	ErrNationalityRequiredForIndividual          = "nationality is required for Individual entity type"
	ErrCountryOfResidenceRequiredForIndividual   = "countryOfResidence is required for Individual entity type"
	ErrCountryOfIncorporationRequiredForBusiness = "countryOfIncorporation is required for Business entity type"
	ErrInvalidEntityType                         = "entityType must be either Individual or Business"
	ErrTokenAlreadyExists                        = "Token already exists"
	ErrCredentialNotFound                        = "Credential not found"
	ErrDocumentNotFound                          = "Document ID not found"
	ErrFileRequired                              = "File is required"
	ErrFailedToOpenFile                          = "Failed to open uploaded file"
	ErrFailedToUploadFile                        = "Failed to upload file"
	ErrFailedToCreateToken                       = "Failed to create token"
	ErrStorageConfigNil                          = "storage config cannot be nil"
	ErrFailedToCreateStorageProvider             = "failed to initialize storage provider"
	ErrFailedToCheckFileExistence                = "failed to check file existence"
	ErrFailedToDeleteFile                        = "failed to delete file"
)

// Error Codes
const (
	ErrCodeBadRequest          = "BAD_REQUEST"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeNotFound            = "RESOURCE_NOT_FOUND"
	ErrCodeConflict            = "CONFLICT"
	ErrCodeInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrCodeValidationFailed    = "VALIDATION_FAILED"
)

// Status Constants
const (
	StatusSuccessful = "successful"
	StatusFailed     = "failed"
	StatusPending    = "Pending"
	StatusActive     = "active"
)

// Entity Constants
const (
	EntityTypeActor  = "ACTOR"
	CredentialTypeVC = "VerifiableCredential"
	TokenStandardVC  = "VC"
)

// Verification Level Constants
const (
	VerificationLevelUnverified = "Tier0_Unverified"
)

// Response Messages
const (
	MsgActorRegisteredSuccessfully     = "Actor registered successfully"
	MsgSignoutSuccessful               = "Signout successful"
	MsgPasswordResetInstructions       = "Password reset instructions have been sent if the email exists"
	MsgEndpointNotFound                = "Endpoint Not Found"
	MsgBadRequest                      = "Bad Request"
	MsgInternalServerError             = "Internal Server Error"
	MsgCredentialSubmittedSuccessfully = "Credential submitted successfully. Verification is in progress."
	MsgOperationSuccessful             = "Operation successful."
	MsgFileUploadedSuccessfully        = "File uploaded successfully"
	MsgUploadedAwaitingVerification    = "Uploaded. Awaiting verification."
)

// HTTP Status Codes
const (
	HTTPStatusOK                  = 200
	HTTPStatusCreated             = 201
	HTTPStatusBadRequest          = 400
	HTTPStatusUnauthorized        = 401
	HTTPStatusNotFound            = 404
	HTTPStatusConflict            = 409
	HTTPStatusInternalServerError = 500
)

// Default Values
const (
	DefaultVerificationLevel = "Tier0_Unverified"
)

// Database Error Constants
const (
	DBErrRecordNotFound = "record not found"
	DBErrDuplicatedKey  = "duplicated key"
)

// Validation Constants
const (
	ValidationTagRequired = "required"
	ValidationTagEmail    = "email"
	ValidationTagMin      = "min=8"
	ValidationTagPassword = "password"
	ValidationTagE164     = "e164"
	ValidationTagOneOf    = "oneof=Individual Business"
	ValidationTagLen2     = "len=2"
	ValidationTagDateTime = "datetime"

	// Field validation constraints
	PasswordMinLength    = 8
	PhoneNumberMinLength = 7
	PhoneNumberMaxLength = 15
	CountryCodeLength    = 2
)

// API Response Field Names
const (
	ResponseFieldID       = "id"
	ResponseFieldVer      = "ver"
	ResponseFieldTs       = "ts"
	ResponseFieldParams   = "params"
	ResponseFieldResponse = "response"
	ResponseFieldStatus   = "status"
	ResponseFieldMessage  = "message"
)

// JWT Response Field Names
const (
	JWTFieldAccessToken  = "accessToken"
	JWTFieldTokenType    = "tokenType"
	JWTFieldExpiresIn    = "expiresIn"
	JWTFieldRefreshToken = "refreshToken"
	JWTFieldScope        = "scope"
)

// Actor Model Field Names
const (
	ActorFieldActorID                = "actor_id"
	ActorFieldDID                    = "did"
	ActorFieldEmail                  = "email"
	ActorFieldFirstName              = "first_name"
	ActorFieldLastName               = "last_name"
	ActorFieldPhoneNumber            = "phone_number"
	ActorFieldMasterPublicKey        = "master_public_key"
	ActorFieldEntityType             = "entity_type"
	ActorFieldVerificationLevel      = "verification_level"
	ActorFieldNationality            = "nationality"
	ActorFieldCountryOfResidence     = "country_of_residence"
	ActorFieldCountryOfIncorporation = "country_of_incorporation"
	ActorFieldCreatedAt              = "created_at"
)

// Identifier Model Field Names
const (
	IdentifierFieldIdentifier = "identifier"
	IdentifierFieldEntityType = "entity_type"
	IdentifierFieldEntityID   = "entity_id"
)

// Actor Profile Field Names
const (
	ActorProfileFieldDID                    = "did"
	ActorProfileFieldUniversalIdentifier    = "universalIdentifier"
	ActorProfileFieldEmail                  = "email"
	ActorProfileFieldFirstName              = "firstName"
	ActorProfileFieldLastName               = "lastName"
	ActorProfileFieldPhoneNumber            = "phoneNumber"
	ActorProfileFieldMasterPublicKey        = "masterPublicKey"
	ActorProfileFieldVerificationLevel      = "verificationLevel"
	ActorProfileFieldEntityType             = "entityType"
	ActorProfileFieldNationality            = "nationality"
	ActorProfileFieldCountryOfResidence     = "countryOfResidence"
	ActorProfileFieldCountryOfIncorporation = "countryOfIncorporation"
)

// Keycloak Constants
const (
	// Grant Types
	KeycloakGrantTypePassword          = "password"
	KeycloakGrantTypeClientCredentials = "client_credentials"

	// Scopes
	KeycloakScopeOpenID = "openid"

	// Clients and Realms
	KeycloakClientAdminCLI = "admin-cli"
	KeycloakRealmMaster    = "master"

	// Content Types
	KeycloakContentTypeJSON = "application/json"
	KeycloakContentTypeForm = "application/x-www-form-urlencoded"

	// Credential Types
	KeycloakCredentialTypePassword = "password"

	// Provider Name
	KeycloakProviderName = "keycloak"

	// API Endpoints
	KeycloakPathAdminUsers              = "/admin/realms/%s/users"
	KeycloakPathAdminUserLogout         = "/admin/realms/%s/users/%s/logout"
	KeycloakPathAdminUserExecuteActions = "/admin/realms/%s/users/%s/execute-actions-email"
	KeycloakPathToken                   = "/realms/%s/protocol/openid-connect/token"
	KeycloakPathMasterToken             = "/realms/master/protocol/openid-connect/token"
	KeycloakPathCerts                   = "/realms/%s/protocol/openid-connect/certs"
	KeycloakPathUserInfo                = "/realms/%s/protocol/openid-connect/userinfo"
	KeycloakPathRealmIssuer             = "%s/realms/%s"
)

// HTTP Header Constants
const (
	HTTPHeaderContentType   = "Content-Type"
	HTTPHeaderAuthorization = "Authorization"
	HTTPHeaderBearer        = "Bearer"
	HTTPHeaderBearerLower   = "bearer"
)

// HTTP Request Parameter Constants
const (
	HTTPParamGrantType    = "grant_type"
	HTTPParamClientID     = "client_id"
	HTTPParamClientSecret = "client_secret"
	HTTPParamUsername     = "username"
	HTTPParamPassword     = "password"
	HTTPParamScope        = "scope"
)

// API Request Defaults
const (
	DefaultRequestID      = "api.request"
	DefaultRequestVersion = "v1"
	DefaultMsgID          = "9cba1b5e-ef63-4a9d-b51d-4b0c03737a32"
)

// Validation Tag Names (for custom error messages)
const (
	ValidationTagNameRequired = "required"
	ValidationTagNameEmail    = "email"
	ValidationTagNameMin      = "min"
	ValidationTagNameMax      = "max"
	ValidationTagNameLen      = "len"
	ValidationTagNameNumber   = "number"
	ValidationTagNamePositive = "positive"
	ValidationTagNameAlphanum = "alphanum"
	ValidationTagNameOneOf    = "oneof"
	ValidationTagNamePassword = "password"
	ValidationTagNameE164     = "e164"
)

// Validation Error Messages
const (
	ValidationMsgRequired = "Field %s must be filled"
	ValidationMsgEmail    = "Invalid email address for field %s"
	ValidationMsgMin      = "Field %s must have a minimum length of %s characters"
	ValidationMsgMax      = "Field %s must have a maximum length of %s characters"
	ValidationMsgLen      = "Field %s must be exactly %s characters long"
	ValidationMsgNumber   = "Field %s must be a number"
	ValidationMsgPositive = "Field %s must be a positive number"
	ValidationMsgAlphanum = "Field %s must contain only alphanumeric characters"
	ValidationMsgOneOf    = "Invalid value for field %s"
	ValidationMsgPassword = "Field %s must contain at least 1 letter and 1 number, minimum length %d"
	ValidationMsgE164     = "Field %s must be a valid phone number in E.164 format (e.g., +1234567890)"
	ValidationMsgDefault  = "Field validation for '%s' failed on the '%s' tag"
)

// Regex Patterns
const (
	RegexDigit      = `[0-9]`
	RegexLetter     = `[a-zA-Z]`
	RegexDigitsOnly = `^[0-9]+$`
)

// Entity Type Constants
const (
	EntityTypeIndividual = "Individual"
	EntityTypeBusiness   = "Business"
)

// Table Names
const (
	TableNameActors            = "actors"
	TableNameIdentifiers       = "identifiers"
	TableNameActorIntegrations = "actor_integrations"
	TableNameTokens            = "tokens"
)

// Database Constants
const (
	DBSSLMode   = "disable"
	DBTimeZone  = "Asia/Shanghai"
	DBDSNFormat = "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s"
)

// Connection Pool Settings
const (
	DBMaxIdleConns    = 10
	DBMaxOpenConns    = 100
	DBConnMaxLifetime = 60 // minutes
)

// Timeout and Cache Duration Constants
const (
	HTTPClientTimeout      = 30 // seconds
	HTTPClientTimeoutShort = 10 // seconds
	KeycloakCacheDuration  = 1  // hours
	StorageURLExpiration   = 24 // hours
)

// Environment Constants
const (
	EnvProduction = "prod"
)

// Environment Variable Names
const (
	EnvAppEnv                = "APP_ENV"
	EnvAppHost               = "APP_HOST"
	EnvAppPort               = "APP_PORT"
	EnvDBHost                = "DB_HOST"
	EnvDBUser                = "DB_USER"
	EnvDBPassword            = "DB_PASSWORD"
	EnvDBName                = "DB_NAME"
	EnvDBPort                = "DB_PORT"
	EnvKeycloakURL           = "KEYCLOAK_URL"
	EnvKeycloakRealm         = "KEYCLOAK_REALM"
	EnvKeycloakClientID      = "KEYCLOAK_CLIENT_ID"
	EnvKeycloakClientSecret  = "KEYCLOAK_CLIENT_SECRET"
	EnvKeycloakAdminUser     = "KEYCLOAK_ADMIN_USER"
	EnvKeycloakAdminPassword = "KEYCLOAK_ADMIN_PASSWORD"
)

// Server Configuration
const (
	ServerHeaderName = "Fiber"
	AppName          = "Fiber API"
)

// Logger Configuration
const (
	LoggerFormat          = "${time} ${method} ${status} ${path} in ${latency}\n"
	LoggerTimeFormat      = "15:04:05.00"
	LoggerTimestampFormat = "15:04:05.000"
)

// Health Check Constants
const (
	HealthStatusUp       = "Up"
	HealthStatusDown     = "Down"
	HealthStatusSuccess  = "success"
	HealthStatusError    = "error"
	HealthServicePostgre = "Postgre"
	HealthServiceMemory  = "Memory"
	HealthCheckCompleted = "Health check completed"
	HealthHeapThreshold  = 300 * 1024 * 1024 // 300 MB
	HealthHeapErrorMsg   = "heap memory usage too high"
)

// API Route Paths
const (
	RouteGroupV1      = "/v1"
	RouteHealthCheck  = "/health-check"
	RouteDocs         = "/docs"
	RouteDocsWildcard = "/*"
)

// Storage Provider Error Messages
const (
	ErrFailedToCreateGCSClient      = "failed to create GCS client"
	ErrFailedToWriteFileToGCS       = "failed to write file to GCS"
	ErrFailedToCloseGCSWriter       = "failed to close GCS writer"
	ErrFailedToGenerateSignedURL    = "failed to generate signed URL"
	ErrFailedToDeleteFileFromGCS    = "failed to delete file from GCS"
	ErrFailedToLoadAWSConfig        = "failed to load AWS config"
	ErrFailedToUploadFileToS3       = "failed to upload file to S3"
	ErrFailedToGeneratePreSignedURL = "failed to generate pre-signed URL"
	ErrFailedToDeleteFileFromS3     = "failed to delete file from S3"
	ErrFailedToCreateMinIOClient    = "failed to create MinIO client"
	ErrFailedToCheckBucketExistence = "failed to check bucket existence"
	ErrFailedToCreateBucket         = "failed to create bucket"
	ErrFailedToUploadFileToMinIO    = "failed to upload file to MinIO"
	ErrFailedToDeleteFileFromMinIO  = "failed to delete file from MinIO"
	ErrSizeMismatch                 = "size mismatch: wrote %d bytes, expected %d"
	ErrNotFound                     = "NotFound"
	ErrNoSuchKey                    = "NoSuchKey"
)

// Storage Provider HTTP Methods
const (
	HTTPMethodGET = "GET"
)

// Storage Provider Signing Schemes
const (
	SigningSchemeV4 = "SigningSchemeV4"
)

// Swagger/API Documentation
const (
	SwaggerTitle    = "Finternet UNITS APIs Documentation"
	SwaggerVersion  = "1.0.0"
	SwaggerLicense  = "MIT"
	SwaggerHost     = "localhost:3000"
	SwaggerBasePath = "/v1"
)

// Phone Number Validation Constants
const (
	PhonePrefixPlus = "+"
	PhoneSpace      = " "
)

// Token expiry example
const (
	TokenExpiryExample = 3600
)
