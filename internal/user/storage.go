package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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

func (u *UserStorage) GetUsers(ctx context.Context, filter *UserFilter, pagination *Pagination) ([]User, error) {
	// Базовый запрос
	query := "SELECT id, name, surname, patronymic, age, sex, nationality FROM users"

	// Добавляем условия фильтрации
	var conditions []string
	var args []any
	argPos := 1

	if filter.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argPos))
		args = append(args, "%"+filter.Name+"%")
		argPos++
	}

	if filter.Surname != "" {
		conditions = append(conditions, fmt.Sprintf("surname ILIKE $%d", argPos))
		args = append(args, "%"+filter.Surname+"%")
		argPos++
	}

	if filter.AgeFrom != -1 {
		conditions = append(conditions, fmt.Sprintf("age >= $%d", argPos))
		args = append(args, filter.AgeFrom)
		argPos++
	}

	if filter.AgeTo != -1 {
		conditions = append(conditions, fmt.Sprintf("age <= $%d", argPos))
		args = append(args, filter.AgeTo)
		argPos++
	}

	if filter.Sex != "" {
		conditions = append(conditions, fmt.Sprintf("sex = $%d", argPos))
		args = append(args, filter.Sex)
		argPos++
	}

	if filter.Nationality != "" {
		conditions = append(conditions, fmt.Sprintf("nationality = $%d", argPos))
		args = append(args, filter.Nationality)
		argPos++
	}

	// Объединяем условия
	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		query += whereClause
	}

	// Добавляем пагинацию
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, pagination.PerPage, (pagination.Page-1)*pagination.PerPage)

	// Выполняем запрос
	rows, err := u.DB.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Surname,
			&user.Patronymic,
			&user.Age,
			&user.Sex,
			&user.Nationality,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
