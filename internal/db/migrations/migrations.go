package migrations

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrationsUp(db *sql.DB) error {
	// Создаем экземпляр драйвера для PostgreSQL
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("невозможно создать драйвер миграции: %w", err)
	}

	// Инициализируем мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("ошибка инициализации мигратора: %w", err)
	}

	// Применяем миграции up
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("невозможно применить up миграции: %w", err)
	}

	log.Println("Миграции up успешно применены.")
	return nil
}

func RunMigrationsDown(db *sql.DB) error {
	// Создаем экземпляр драйвера для PostgreSQL
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("невозможно создать драйвер миграции: %w", err)
	}

	// Инициализируем мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("ошибка инициализации мигратора: %w", err)
	}

	// Применяем миграции down
	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("невозможно применить down миграции: %w", err)
	}

	log.Println("Миграции down успешно применены.")
	return nil
}
