package middleware

import (
	"app/src/config"
	"app/src/constants"
	"app/src/model"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// JWKSKey represents a single key from JWKS endpoint
type JWKSKey struct {
	Kid string   `json:"kid"`
	Kty string   `json:"kty"`
	Alg string   `json:"alg"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// JWKSResponse represents the JWKS endpoint response
type JWKSResponse struct {
	Keys []JWKSKey `json:"keys"`
}

// AuthClaims represents the JWT claims from the auth provider
type AuthClaims struct {
	Sub               string `json:"sub"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	jwt.RegisteredClaims
}

// AuthJWTValidator handles JWT validation with the auth provider
type AuthJWTValidator struct {
	log           *logrus.Logger
	baseURL       string
	realm         string
	clientID      string
	publicKeys    map[string]*rsa.PublicKey
	keysMutex     sync.RWMutex
	lastFetch     time.Time
	cacheDuration time.Duration
	db            *gorm.DB
	httpClient    *http.Client
}

// NewAuthJWTValidator creates a new auth JWT validator
func NewAuthJWTValidator(cfg *config.Config, log *logrus.Logger, db *gorm.DB) (*AuthJWTValidator, error) {
	if cfg.AuthURL == "" || cfg.AuthRealm == "" || cfg.AuthClientID == "" {
		return nil, fmt.Errorf("auth configuration incomplete: URL, realm, and client ID are required")
	}

	validator := &AuthJWTValidator{
		log:           log,
		baseURL:       cfg.AuthURL,
		realm:         cfg.AuthRealm,
		clientID:      cfg.AuthClientID,
		publicKeys:    make(map[string]*rsa.PublicKey),
		cacheDuration: constants.KeycloakCacheDuration * time.Hour,
		db:            db,
		httpClient: &http.Client{
			Timeout: constants.HTTPClientTimeoutShort * time.Second,
		},
	}

	// Pre-fetch keys on initialization
	if err := validator.fetchPublicKeys(); err != nil {
		return nil, fmt.Errorf("failed to fetch auth provider public keys: %w", err)
	}

	return validator, nil
}

// AuthMiddleware handles authentication middleware
type AuthMiddleware struct {
	validator *AuthJWTValidator
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(validator *AuthJWTValidator) *AuthMiddleware {
	return &AuthMiddleware{validator: validator}
}

// Authenticate returns a fiber.Handler that validates JWT tokens
func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := m.extractBearerToken(c)
		if err != nil {
			return err
		}

		claims, err := m.validator.ValidateToken(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, constants.ErrUnauthorized)
		}

		if err := m.validator.ValidateSessionActive(token); err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, constants.ErrUnauthorized)
		}

		actorID, err := m.validator.FindActorByUsername(c.Context(), claims.PreferredUsername, claims.Email)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, constants.ErrUnauthorized)
		}

		c.Locals("actorID", actorID)
		c.Locals("auth_claims", claims)

		return c.Next()
	}
}

// extractBearerToken extracts and validates the Bearer token from the Authorization header
func (m *AuthMiddleware) extractBearerToken(c *fiber.Ctx) (string, error) {
	authHeader := c.Get(constants.HTTPHeaderAuthorization)
	if authHeader == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, constants.ErrUnauthorized)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != constants.HTTPHeaderBearerLower {
		return "", fiber.NewError(fiber.StatusUnauthorized, constants.ErrUnauthorized)
	}

	return parts[1], nil
}

// ValidateToken validates a JWT token using RSA signature verification
func (v *AuthJWTValidator) ValidateToken(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("missing kid in JWT header")
		}

		return v.GetPublicKey(kid)
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Validate issuer
	expectedIssuer := fmt.Sprintf(constants.KeycloakPathRealmIssuer, v.baseURL, v.realm)
	if claims.Issuer != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer: expected %s, got %s", expectedIssuer, claims.Issuer)
	}

	return claims, nil
}

// GetPublicKey retrieves the public key for the given kid
func (v *AuthJWTValidator) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	v.keysMutex.RLock()
	key, exists := v.publicKeys[kid]
	shouldRefresh := time.Since(v.lastFetch) > v.cacheDuration
	v.keysMutex.RUnlock()

	if exists && !shouldRefresh {
		return key, nil
	}

	if err := v.fetchPublicKeys(); err != nil {
		return nil, err
	}

	v.keysMutex.RLock()
	defer v.keysMutex.RUnlock()

	key, exists = v.publicKeys[kid]
	if !exists {
		return nil, fmt.Errorf("public key not found for kid: %s", kid)
	}

	return key, nil
}

// fetchPublicKeys fetches public keys from the auth provider JWKS endpoint
func (v *AuthJWTValidator) fetchPublicKeys() error {
	jwksURL := fmt.Sprintf("%s"+constants.KeycloakPathCerts, v.baseURL, v.realm)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch JWKS: HTTP %d", resp.StatusCode)
	}

	var jwks JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	v.keysMutex.Lock()
	defer v.keysMutex.Unlock()

	for _, key := range jwks.Keys {
		if key.Kty == "RSA" {
			if publicKey, err := v.parseRSAPublicKey(key); err == nil {
				v.publicKeys[key.Kid] = publicKey
			}
		}
	}

	v.lastFetch = time.Now()
	return nil
}

// parseRSAPublicKey converts JWKS key to RSA public key
func (v *AuthJWTValidator) parseRSAPublicKey(key JWKSKey) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode N: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode E: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := 0
	for _, b := range eBytes {
		e = e*256 + int(b)
	}

	return &rsa.PublicKey{N: n, E: e}, nil
}

// ValidateSessionActive checks if the token's session is still active
func (v *AuthJWTValidator) ValidateSessionActive(token string) error {
	userinfoURL := fmt.Sprintf("%s"+constants.KeycloakPathUserInfo, v.baseURL, v.realm)

	req, err := http.NewRequest("GET", userinfoURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set(constants.HTTPHeaderAuthorization, constants.HTTPHeaderBearer+" "+token)

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid session")
	}

	return nil
}

// FindActorByUsername finds an actor ID by username or email
func (v *AuthJWTValidator) FindActorByUsername(ctx context.Context, username, email string) (uuid.UUID, error) {
	var actor model.Actor

	// Try by universal identifier first
	var identifier model.Identifier
	if v.db.WithContext(ctx).Where("identifier = ?", username).First(&identifier).Error == nil {
		if v.db.WithContext(ctx).Where("actor_id = ?", identifier.EntityID).First(&actor).Error == nil {
			return actor.ActorID, nil
		}
	}

	// Try by email
	if v.db.WithContext(ctx).Where("email = ?", email).First(&actor).Error == nil {
		return actor.ActorID, nil
	}

	return uuid.Nil, fmt.Errorf("actor not found")
}
