package mysql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(dsn string) (*UserRepository, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &UserRepository{pool: pool}, nil
}

func (r *UserRepository) CreateEmp(empType, username, passwordHash string) error {
	_, err := r.pool.Exec(context.Background(),
		"INSERT INTO employees (emp_type, username, password_hash) VALUES ($1, $2, $3)",
		empType, username, passwordHash,
	)
	return err
}

func (r *UserRepository) GetEmpByEmpName(username string) (int, string, error) {
	var empID int
	var passwordHash string
	err := r.pool.QueryRow(context.Background(),
		"SELECT id, password_hash FROM employees WHERE username = $1", username,
	).Scan(&empID, &passwordHash)

	if err != nil {
		return 0, "", errors.New("employee not found")
	}
	return empID, passwordHash, nil
}
