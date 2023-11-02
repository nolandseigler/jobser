package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

var (
	genSaltSQL = "gen_salt('bf', 8)"
)

func New(ctx context.Context, config Config) (*DB, error) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	pool, err := pgxpool.New(
		ctx,
		fmt.Sprintf(
			"%s://%s:%s@%s:%d/%s",
			config.Protocol,
			config.Username,
			config.Password,
			config.Hostname,
			config.Port,
			config.DatabaseName,
		),
	)
	if err != nil {
		return nil, err
	}
	return &DB{
		pool: pool,
	}, nil
}

func (d *DB) Close(ctx context.Context) {
	d.pool.Close()
}

func (d *DB) GetUserAccount(ctx context.Context, username string) (*UserAccount, error) {
	var users []*UserAccount
	err := pgxscan.Select(
		ctx,
		d.pool,
		&users,
		`SELECT username FROM auth.user_account WHERE username = $1`,
		username,
	)
	if err != nil {
		return nil, err
	}

	usersLen := len(users)

	if usersLen == 0 {
		return nil, errors.New("user not found")
	}
	if usersLen > 1 {
		return nil, errors.New("more than one user found")
	}

	return users[0], nil
}

func (d *DB) CreateUserAccount(ctx context.Context, username string, password string) (*UserAccount, error) {

	passwordLen := len(password)
	if passwordLen < 12 || passwordLen > 60 {
		return nil, fmt.Errorf("password length must be >= 12 and <= 60; password length: %d", passwordLen)
	}
	tx, err := d.pool.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := tx.Exec(
		ctx,
		fmt.Sprintf(
			"insert into auth.user_account (username, password) values ($1, crypt($2, %s));",
			genSaltSQL,
		),
		username,
		password,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx) ;err != nil {
		return nil, err
	}


	return &UserAccount{
		Username: username,
	}, nil
}

func (d *DB) IsUserAccountPassword(ctx context.Context, username string, password string) (bool, error) {
	var users []*UserAccount
	err := pgxscan.Select(
		ctx,
		d.pool,
		&users,
		fmt.Sprintf(
			`SELECT username FROM auth.user_account WHERE username = $1 and password = crypt($2, %s)`,
			genSaltSQL,
		),
		username,
	)
	if err != nil {
		return false, err
	}

	usersLen := len(users)

	if usersLen == 0 {
		return false, errors.New("user not found")
	}
	if usersLen > 1 {
		return false, errors.New("more than one user found")
	}

	return true, nil
}
