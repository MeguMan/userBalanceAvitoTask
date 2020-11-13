package postgres_store

import (
	"github.com/jackc/pgx/v4"
)

type Store struct {
	conn *pgx.Conn
	UserRepository *UserRepository
}

func New(conn *pgx.Conn) *Store {
	return &Store{
		conn: conn,
	}
}

func (s *Store) User() *UserRepository {
	if s.UserRepository == nil {
		s.UserRepository = &UserRepository{
			store: s,
		}
	}

	return s.UserRepository
}