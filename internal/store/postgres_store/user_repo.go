package postgres_store

import (
	"context"
	"errors"
	"github.com/MeguMan/userBalanceAvitoTask/internal/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) IncreaseBalance(u model.User) error {
	currentBalance, err := r.GetBalanceById(u.ID)
	if err != nil {
		return err
	}

	_, err = r.store.conn.Exec(context.Background(),"UPDATE users SET balance=$1 where user_id=$2",
		currentBalance + u.Balance, u.ID)

	return err
}

func (r *UserRepository) ReduceBalance(u model.User) error {
	currentBalance, err := r.GetBalanceById(u.ID)
	if err != nil {
		return err
	}
	if u.Balance > currentBalance {
		return errors.New("balance can't be less than 0")
	}

	_, err = r.store.conn.Exec(context.Background(),"UPDATE users SET balance=$1 where user_id=$2",
		currentBalance - u.Balance, u.ID)

	return err
}

func (r *UserRepository) GetBalanceById(id int) (float64, error) {
	var balance float64
	err := r.store.conn.QueryRow(context.Background(),"SELECT balance FROM users WHERE user_id=$1",
		id).Scan(&balance)

	return balance, err
}

func (r *UserRepository) TransferBalance(sender, receiver model.User) error {
	err := r.ReduceBalance(sender)
	if err != nil {
		 return err
	}
	err = r.IncreaseBalance(receiver)
	return err
}
