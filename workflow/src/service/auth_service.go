package service

import (
	"app/src/config"
	"app/src/constants"
	"app/src/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// AuthService provides authentication and user management operations
type AuthService interface {
	CreateUser(actor *model.Actor, universalIdentifier, password string) (string, error)
	Login(username, password string) (*AuthTokenResponse, error)
	Logout(userID string) error
	ExecuteActionsEmail(userID string, actions []string) error
}

// authService implements AuthService using OIDC-compliant auth provider
type authService struct {
	log          *logrus.Logger
	baseURL      string
	realm        string
	clientID     string
	clientSecret string
	adminUser    string
	adminPass    string
	httpClient   *http.Client
}

// AuthUserRequest represents a user creation request
type AuthUserRequest struct {
	Username      string           `json:"username"`
	Email         string           `json:"email"`
	FirstName     string           `json:"firstName"`
	LastName      string           `json:"lastName"`
	Enabled       bool             `json:"enabled"`
	EmailVerified bool             `json:"emailVerified"`
	Credentials   []AuthCredential `json:"credentials"`
}

// AuthCredential represents a user credential
type AuthCredential struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

// AuthTokenResponse represents an authentication token response
type AuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// AuthErrorResponse represents an error response from the auth provider
type AuthErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

// NewAuthService creates a new auth service instance
func NewAuthService(cfg *config.Config, log *logrus.Logger) AuthService {
	return &authService{
		log:          log,
		baseURL:      cfg.AuthURL,
		realm:        cfg.AuthRealm,
		clientID:     cfg.AuthClientID,
		clientSecret: cfg.AuthSecret,
		adminUser:    cfg.AuthAdminUser,
		adminPass:    cfg.AuthAdminPassword,
		httpClient: &http.Client{
			Timeout: constants.HTTPClientTimeout * time.Second,
		},
	}
}

func (s *authService) CreateUser(actor *model.Actor, universalIdentifier, password string) (string, error) {
	adminToken, err := s.getAdminToken()
	if err != nil {
		return "", fmt.Errorf("failed to authenticate with auth provider: %w", err)
	}

	userReq := AuthUserRequest{
		Username:      universalIdentifier,
		Email:         actor.Email,
		FirstName:     actor.FirstName,
		LastName:      actor.LastName,
		Enabled:       true,
		EmailVerified: true,
		Credentials: []AuthCredential{
			{Type: constants.KeycloakCredentialTypePassword, Value: password, Temporary: false},
		},
	}

	jsonData, err := json.Marshal(userReq)
	if err != nil {
		return "", fmt.Errorf("failed to prepare user creation request: %w", err)
	}

	apiURL := fmt.Sprintf("%s"+constants.KeycloakPathAdminUsers, s.baseURL, s.realm)
	resp, err := s.doRequest("POST", apiURL, adminToken, bytes.NewBuffer(jsonData), constants.KeycloakContentTypeJSON)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Handle user already exists
	if resp.StatusCode == http.StatusConflict {
		s.log.Infof("User already exists for: %s, fetching user ID", universalIdentifier)
		return s.getUserIDByUsername(adminToken, universalIdentifier)
	}

	// Handle creation errors
	if resp.StatusCode != http.StatusCreated {
		return "", s.handleErrorResponse(resp, "failed to create user")
	}

	// Extract user ID from Location header
	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("failed to get user ID from response")
	}

	parts := strings.Split(location, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid Location header format")
	}
	userID := parts[len(parts)-1]

	s.log.Infof("Successfully created user for: %s with ID: %s", universalIdentifier, userID)
	return userID, nil
}

func (s *authService) Login(username, password string) (*AuthTokenResponse, error) {
	s.log.Infof("Login attempt for user: %s", username)

	tokenURL := fmt.Sprintf("%s"+constants.KeycloakPathToken, s.baseURL, s.realm)

	data := url.Values{}
	data.Set(constants.HTTPParamGrantType, constants.KeycloakGrantTypePassword)
	data.Set(constants.HTTPParamClientID, s.clientID)
	data.Set(constants.HTTPParamClientSecret, s.clientSecret)
	data.Set(constants.HTTPParamUsername, username)
	data.Set(constants.HTTPParamPassword, password)
	data.Set(constants.HTTPParamScope, constants.KeycloakScopeOpenID)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set(constants.HTTPHeaderContentType, constants.KeycloakContentTypeForm)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.log.Errorf("Authentication failed for user: %s", username)
		return nil, fiber.NewError(fiber.StatusUnauthorized, constants.ErrInvalidCredentials)
	}

	var tokenResp AuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode authentication response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		s.log.Errorf("Received empty access token for user: %s", username)
		return nil, fmt.Errorf("received empty access token")
	}

	s.log.Infof("Successfully authenticated user: %s", username)
	return &tokenResp, nil
}

func (s *authService) Logout(userID string) error {
	s.log.Infof("Initiating logout for user: %s", userID)

	adminToken, err := s.getAdminToken()
	if err != nil {
		return fmt.Errorf("failed to authenticate with auth provider: %w", err)
	}

	logoutURL := fmt.Sprintf("%s"+constants.KeycloakPathAdminUserLogout, s.baseURL, s.realm, userID)
	resp, err := s.doRequest("POST", logoutURL, adminToken, nil, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return s.handleErrorResponse(resp, "failed to logout user")
	}

	s.log.Infof("Successfully logged out user: %s", userID)
	return nil
}

func (s *authService) ExecuteActionsEmail(userID string, actions []string) error {
	s.log.Infof("Executing actions email for user: %s, actions: %v", userID, actions)

	adminToken, err := s.getAdminToken()
	if err != nil {
		return fmt.Errorf("failed to authenticate with auth provider: %w", err)
	}

	executeActionsURL := fmt.Sprintf("%s"+constants.KeycloakPathAdminUserExecuteActions, s.baseURL, s.realm, userID)

	jsonData, err := json.Marshal(actions)
	if err != nil {
		return fmt.Errorf("failed to prepare execute actions request: %w", err)
	}

	resp, err := s.doRequest("PUT", executeActionsURL, adminToken, bytes.NewBuffer(jsonData), constants.KeycloakContentTypeJSON)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return s.handleErrorResponse(resp, "failed to execute actions email")
	}

	s.log.Infof("Successfully sent actions email to user: %s", userID)
	return nil
}

// getAdminToken retrieves an admin token for API operations
func (s *authService) getAdminToken() (string, error) {
	tokenURL := s.baseURL + constants.KeycloakPathMasterToken

	data := url.Values{}
	data.Set(constants.HTTPParamGrantType, constants.KeycloakGrantTypePassword)
	data.Set(constants.HTTPParamClientID, constants.KeycloakClientAdminCLI)
	data.Set(constants.HTTPParamUsername, s.adminUser)
	data.Set(constants.HTTPParamPassword, s.adminPass)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create admin token request: %w", err)
	}

	req.Header.Set(constants.HTTPHeaderContentType, constants.KeycloakContentTypeForm)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get admin token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("admin token request failed: HTTP %d", resp.StatusCode)
	}

	var tokenResp AuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode admin token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

// getUserIDByUsername queries the auth provider to get a user's ID by username
func (s *authService) getUserIDByUsername(adminToken, username string) (string, error) {
	apiURL := fmt.Sprintf("%s"+constants.KeycloakPathAdminUsers+"?username=%s&exact=true",
		s.baseURL, s.realm, url.QueryEscape(username))

	resp, err := s.doRequest("GET", apiURL, adminToken, nil, "")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to query user: HTTP %d", resp.StatusCode)
	}

	var users []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return "", fmt.Errorf("failed to decode user query response: %w", err)
	}

	if len(users) == 0 {
		return "", fmt.Errorf("user not found: %s", username)
	}

	userID, ok := users[0]["id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract user ID from response")
	}

	s.log.Infof("Found user ID %s for username: %s", userID, username)
	return userID, nil
}

// doRequest is a helper method to perform HTTP requests with authorization
func (s *authService) doRequest(method, url, token string, body *bytes.Buffer, contentType string) (*http.Response, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, body)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		s.log.Errorf("Failed to create %s request: %v", method, err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if token != "" {
		req.Header.Set(constants.HTTPHeaderAuthorization, constants.HTTPHeaderBearer+" "+token)
	}

	if contentType != "" {
		req.Header.Set(constants.HTTPHeaderContentType, contentType)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.log.Errorf("Failed to execute %s request: %v", method, err)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// handleErrorResponse processes and returns an error from an HTTP response
func (s *authService) handleErrorResponse(resp *http.Response, baseMsg string) error {
	var errorResp AuthErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		s.log.Errorf("%s: HTTP %d", baseMsg, resp.StatusCode)
		return fmt.Errorf("%s: HTTP %d", baseMsg, resp.StatusCode)
	}
	s.log.Errorf("%s: %s - %s", baseMsg, errorResp.Error, errorResp.Description)
	return fmt.Errorf("%s: %s - %s", baseMsg, errorResp.Error, errorResp.Description)
}
