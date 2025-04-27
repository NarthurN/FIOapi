package postgresdb

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(DBpath string) (*sql.DB, error) {
	op := "internal/db/postgredb/postgredb.go.New"

	db, err := sql.Open("pgx", DBpath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
