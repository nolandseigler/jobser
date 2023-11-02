package handlers

import (
	"context"
	"github.com/nolandseigler/wordser/wordserweb/internal/storage/postgres"
)

type DBer interface {
	CreateUserAccount(ctx context.Context, username string, password string) (*postgres.UserAccount, error)
	GetUserAccount(ctx context.Context, username string) (*postgres.UserAccount, error)
}
