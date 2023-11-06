package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	pubKey       *rsa.PublicKey
	signingKey   *rsa.PrivateKey
	kvStore      KeyValStorer
	userVerifier UserVerifier
}

func New(ctx context.Context, config Config, kvStore KeyValStorer, userVerifier UserVerifier) (*Auth, error) {
	signBytes, err := os.ReadFile(config.PrivKeyPath)
	if err != nil {
		return nil, err
	}

	verifyBytes, err := os.ReadFile(config.PubKeyPath)
	if err != nil {
		return nil, err
	}

	signingKey, err := jwtlib.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}

	pubKey, err := jwtlib.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}

	return &Auth{
		pubKey:       pubKey,
		signingKey:   signingKey,
		kvStore:      kvStore,
		userVerifier: userVerifier,
	}, nil
}

type UserContext struct {
	Username string `json:"username"`
}

type JWTClaims struct {
	UserContext UserContext `json:"user_context"`
	jwtlib.RegisteredClaims
}

// destroySession -> KeyValStorer.Delete(jti)
func (a *Auth) destroySesion(ctx context.Context, jti uuid.UUID) error {
	err := a.kvStore.Delete(jti.String())
	// TODO: maybe handle me
	if err != nil {
		return err
	}
	return nil
}

// storeSession -> KeyValStorer.Insert(jti)
func (a *Auth) storeSession(ctx context.Context, jti uuid.UUID) error {
	err := a.kvStore.Insert(jti.String(), "")
	// TODO: maybe handle me
	if err != nil {
		return err
	}
	return nil
}

// validateSession -> KeyValStorer.Get(jti)
func (a *Auth) validateSession(ctx context.Context, jti uuid.UUID) error {
	if _, ok := a.kvStore.Get(jti.String()); !ok {
		return fmt.Errorf("no valid session")
	}
	return nil
}

// mintJWt
func (a *Auth) mintJWT(ctx context.Context, userCtx UserContext) (string, uuid.UUID, error) {
	jti := uuid.New()
	jwt, err := jwtlib.NewWithClaims(
		jwtlib.SigningMethodRS512,
		JWTClaims{
			userCtx,
			jwtlib.RegisteredClaims{
				// A usual scenario is to set the expiration time relative to the current time
				ExpiresAt: jwtlib.NewNumericDate(time.Now().UTC().Add(24 * time.Hour)),
				IssuedAt:  jwtlib.NewNumericDate(time.Now().UTC()),
				NotBefore: jwtlib.NewNumericDate(time.Now().UTC()),
				Issuer:    "wordserweb",
				Subject:   userCtx.Username,
				ID:        jti.String(),
				Audience:  []string{"wordserweb"},
			},
		},
	).SignedString(a.signingKey)

	if err != nil {
		return "", uuid.Nil, nil
	}

	return jwt, jti, nil
}

// newJWT -> mintJwt, and storeSession
func (a *Auth) newJWT(ctx context.Context, username string) (string, uuid.UUID, error) {
	// TODO: Lookup user context
	jwt, jti, err := a.mintJWT(
		ctx,
		UserContext{
			Username: username,
		},
	)
	if err != nil {
		return "", uuid.Nil, nil
	}

	if err := a.storeSession(ctx, jti); err != nil {
		return "", uuid.Nil, nil
	}

	return jwt, jti, nil
}

// Login -> Use db to check user & pass then NewJWT
func (a *Auth) Login(ctx context.Context, username string, password string) (string, error) {
	verified, err := a.userVerifier.IsUserAccountPassword(ctx, username, password)
	if err != nil {
		return "", err
	}

	if !verified {
		return "", fmt.Errorf("invalid username/password or user does not exist")
	}

	jwt, _, err := a.newJWT(ctx, username)

	if err != nil {
		return "", err
	}

	return jwt, nil
}

// refreshJWT -> part of validateJWT. if validJWT will expire in `n` minutes then NewJWT
func (a *Auth) refreshJWT(ctx context.Context, jti uuid.UUID, username string) (string, uuid.UUID, error) {
	if err := a.destroySesion(ctx, jti); err != nil {
		return "", uuid.Nil, err
	}
	jwt, jti, err := a.newJWT(ctx, username)

	if err != nil {
		return "", uuid.Nil, err
	}

	return jwt, jti, nil
}

// validateJWT -> does the work but isnt the public function. takes an argument refresh: bool so we can use this in logout without refresh.
func (a *Auth) validateJWT(ctx context.Context, jwt string, refresh bool) (string, uuid.UUID, error) {
	token, err := jwtlib.ParseWithClaims(jwt, &JWTClaims{}, func(token *jwtlib.Token) (interface{}, error) {
		return a.pubKey, nil
	})
	if err != nil {
		return "", uuid.Nil, err
	}

	if !token.Valid {
		return "", uuid.Nil, fmt.Errorf("invalid token;")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return "", uuid.Nil, fmt.Errorf("invalid token claims;")
	}

	expiry, err := claims.RegisteredClaims.GetExpirationTime()
	if err != nil {
		return "", uuid.Nil, err
	}

	returnJti, err := uuid.Parse(claims.RegisteredClaims.ID)
	if err != nil {
		return "", uuid.Nil, err
	}

	if err := a.validateSession(ctx, returnJti); err != nil {
		return "", uuid.Nil, err
	}

	returnJwt := jwt

	if refresh && expiry.Time.After(time.Now().UTC().Add(5*time.Minute)) {
		subj, err := claims.RegisteredClaims.GetSubject()
		if err != nil {
			return "", uuid.Nil, err
		}

		returnJwt, returnJti, err = a.refreshJWT(ctx, returnJti, subj)

		if err != nil {
			return "", uuid.Nil, err
		}
		// success continue to return
	}

	return returnJwt, returnJti, nil
}

// Logout -> this is not behind ValidateJWT middleware func so we call validateJWT with refresh = false
// no jwt no logout.
func (a *Auth) Logout(ctx context.Context, jwt string) error {
	_, _, err := a.validateJWT(ctx, jwt, false)
	if err != nil {
		return err
	}
	return nil
}

// ValidateJWT -> includes using redis storage to check session. Validate and call refresh if close to expiry
func (a *Auth) ValidateJWTMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		if c.Path() == "/metrics" ||
			c.Path() == "/signup" ||
			c.Path() == "/login" ||
			strings.HasPrefix(c.Path(), "/static") {
			return next(c)
		}

		sessionCookie, err := c.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				return echo.NewHTTPError(http.StatusUnauthorized, "no credentials")
			}
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}
		jwt, _, err := a.validateJWT(c.Request().Context(), sessionCookie.Value, true)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}

		// TODO: If I end up doing perms per endpoint then put a sync.Map in *Auth that maps endpoint to all []Perms
		// or something.

		if sessionCookie.Value == jwt {
			return next(c)
		}
		// new jwt replace cookie.
		c.SetCookie(&http.Cookie{
			Name:    "session_token",
			Value:   jwt,
			Expires: time.Now().UTC().Add(60 * time.Minute),
		})

		return next(c)
	}
}
