package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Auth struct {
	pubKey     *rsa.PublicKey
	signingKey *rsa.PrivateKey
	kvStore    KeyValStorer
	userVerifier UserVerifier
}

func New(config Config, kvStore KeyValStorer, userVerifier UserVerifier) (*Auth, error) {
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
		pubKey:     pubKey,
		signingKey: signingKey,
		kvStore:    kvStore,
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

// isSessionValid -> KeyValStorer.Get(jti)
func (a *Auth) isSessionValid(ctx context.Context, jti uuid.UUID) error {
	if 	_, ok := a.kvStore.Get(jti.String()); !ok {
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
	err := a.userVerifier.Verify(username, password)
	if err != nil {
		return "", err
	}

	jwt, _, err := a.newJWT(ctx, username)

	if err != nil {
		return "", err
	}

	return jwt, nil
}

// refreshJWT -> part of validateJWT. if validJWT will expire in `n` minutes then NewJWT
func (a *Auth) refreshJWT(ctx context.Context, jti uuid.UUID, username string) (string, uuid.UUID, error) {
	if err := a.kvStore.Delete(jti.String()); err != nil {
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
	returnJwt := jwt

	if refresh && expiry.Time.After(time.Now().UTC().Add(5 * time.Minute)) {
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
func (a *Auth) ValidateJWT(ctx context.Context, jwt string) (string, error) {
	jwt, _, err := a.validateJWT(ctx, jwt, true)
	if err != nil {
		return "", err
	}
	return jwt, err
}
