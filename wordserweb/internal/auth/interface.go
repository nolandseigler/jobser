package auth

import "context"

type KeyValStorer interface {
	Insert(key string, value string) error
	Delete(key string) error
	Get(key string) (string, bool)
}

type UserVerifier interface {
	IsUserAccountPassword(ctx context.Context, username string, password string) (bool, error)
}
