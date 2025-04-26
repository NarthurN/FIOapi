package migrations

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrationsUp(db *sql.DB) error {
	op := "internal/db/migrations/migrations.go.RunMigrationsUp"
	// Создаем экземпляр драйвера для PostgreSQL
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Инициализируем мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Применяем миграции up
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func RunMigrationsDown(db *sql.DB) error {
	op := "internal/db/migrations/migrations.go.RunMigrationsDown"
	// Создаем экземпляр драйвера для PostgreSQL
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Инициализируем мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Применяем миграции down
	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
