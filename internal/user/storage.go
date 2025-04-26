package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NarthurN/FIOapi/internal/db/postgresdb"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Storage - интерфейс для работы с пользователями в БД.
type Storage interface {
	Create(ctx context.Context, user *User) error
	// GetByID(ctx context.Context, id int) (*User, error)
	// GetByEmail(ctx context.Context, email string) (*User, error)
	// Update(ctx context.Context, user *User) error
	// Delete(ctx context.Context, id int) error
}

// PostgresStorage реализует Storage для PostgreSQL.
type UserStorage struct {
	DB *sql.DB
}

func NewStorage() (*UserStorage, error) {
	op := "internal/user/storage.go.NewStorage"
	db, err := postgresdb.New()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &UserStorage{DB: db}, nil
}

// Добавляет пользователя и возвращает id добавленного пользователя.
func (u *UserStorage) Create(ctx context.Context, user *User) error {
	return nil
}
