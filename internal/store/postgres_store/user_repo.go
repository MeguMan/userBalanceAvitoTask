package postgres_store

import (
	"context"
	"errors"
	"github.com/MeguMan/userBalanceAvitoTask/internal/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) AddBalance(u model.User) error {
	currentBalance, err := r.GetBalanceById(u)

	_, err = r.store.conn.Exec(context.Background(),"UPDATE users SET balance=$1 where user_id=$2",
		currentBalance + u.Balance, u.ID)

	return err
}

func (r *UserRepository) ReduceBalance(u model.User) error {
	currentBalance, err := r.GetBalanceById(u)
	if u.Balance > currentBalance {
		return errors.New("balance can't be less than 0")
	}

	_, err = r.store.conn.Exec(context.Background(),"UPDATE users SET balance=$1 where user_id=$2",
		currentBalance - u.Balance, u.ID)

	return err
}

func (r *UserRepository) GetBalanceById(u model.User) (int, error) {
	var balance int
	err := r.store.conn.QueryRow(context.Background(),"SELECT balance FROM users WHERE user_id=$1",
		u.ID).Scan(&balance)

	return balance, err
}
