package postgres_store

import (
	"context"
	"github.com/MeguMan/userBalanceAvitoTask/internal/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) (model.User, error) {
	err := r.store.conn.QueryRow(context.Background(),"INSERT INTO users (balance) VALUES ($1) returning user_id",
		u.Balance).Scan(&u.ID)

	return *u, err
}
