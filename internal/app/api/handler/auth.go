package handler

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	cachedJWKS jwk.Set
	mu         = &sync.RWMutex{}
	lastFetch  time.Time
)

const jwksCacheDuration = 15 * time.Minute

type Auth struct {
	skipPaths []string
}

func NewAuth(skipPaths []string) *Auth {
	return &Auth{
		skipPaths: skipPaths,
	}
}

func (a *Auth) ClerkJWTMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:Authorization",
		AuthScheme: "Bearer",
		Skipper:    a.Skipper,
		Validator:  a.Validator,
	})
}

func (a *Auth) Skipper(c echo.Context) bool {
	for _, p := range a.skipPaths {
		if ok := strings.HasPrefix(c.Path(), p); ok {
			return true
		}
	}
	return false
}

func (a *Auth) Validator(tokenStr string, c echo.Context) (bool, error) {
	jwks, err := getCachedJWKS(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed to get cached jwks: %w", err)
	}

	token, err := jwt.ParseString(tokenStr, jwt.WithKeySet(jwks), jwt.WithValidate(true))
	if err != nil {
		return false, fmt.Errorf("failed to parse tokenStr: %w", err)
	}

	sub := token.Subject()
	c.Set(userIDKey, sub)

	return true, nil
}

func getCachedJWKS(ctx context.Context) (jwk.Set, error) {
	mu.Lock()
	defer mu.Unlock()

	if cachedJWKS != nil && time.Since(lastFetch) < jwksCacheDuration {
		return cachedJWKS, nil
	}

	jwksURL := os.Getenv("CLERK_JWKS_URL")
	if jwksURL == "" {
		// dev
		jwksURL = "https://precious-ghoul-88.clerk.accounts.dev/.well-known/jwks.json"
	}
	newJWKS, err := jwk.Fetch(ctx, jwksURL)
	if err != nil {
		return nil, err
	}

	cachedJWKS = newJWKS
	lastFetch = time.Now()
	return cachedJWKS, nil
}
