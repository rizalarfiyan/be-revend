package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizalarfiyan/be-revend/internal/sql"
)

type userRepository struct {
	db    *pgxpool.Pool
	query *sql.Queries
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db:    db,
		query: sql.New(db),
	}
}
