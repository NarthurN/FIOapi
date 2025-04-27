package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NarthurN/FIOapi/internal/db/postgresdb"
)

// UserStorage реализует Storage для PostgreSQL.
type UserStorage struct {
	DB *sql.DB
}

func NewStorage(DBpath string) (*UserStorage, error) {
	op := "internal/user/storage.go.NewStorage"
	db, err := postgresdb.New(DBpath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &UserStorage{DB: db}, nil
}

// Добавляет пользователя и возвращает id добавленного пользователя.
func (u *UserStorage) Create(ctx context.Context, user *User) (int, error) {
	var id int
	stmt := `INSERT INTO users (name, surname, patronymic, age, sex, nationality)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	row := u.DB.QueryRowContext(ctx, stmt,
		user.Name,
		user.Surname,
		user.Patronymic,
		user.Age,
		user.Sex,
		user.Nationality,
	)
	return id, row.Scan(&id)
}
