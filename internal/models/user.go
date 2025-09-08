package models

import (
	"context"
	"database/sql"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
}

type UserStore struct{ DB *sql.DB }

func (s *UserStore) Create(ctx context.Context, email, passwordHash string) (User, error) {
	var u User
	q := `insert into users(id,email,password_hash) values (gen_random_uuid(),$1,$2) returning id,email,password_hash`
	if err := s.DB.QueryRowContext(ctx, q, email, passwordHash).Scan(&u.ID, &u.Email, &u.PasswordHash); err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *UserStore) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	q := `select id,email,password_hash from users where email=$1`
	if err := s.DB.QueryRowContext(ctx, q, email).Scan(&u.ID, &u.Email, &u.PasswordHash); err != nil {
		return User{}, err
	}
	return u, nil
}
