package postgresdb

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func New() (*sql.DB, error) {
	op := "internal/db/postgredb/postgredb.go.New"
	connStr := "user=postgres password=root host=localhost port=5432 database=postgres sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
