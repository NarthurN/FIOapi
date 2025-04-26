package user

import (
	//"context"
	"database/sql"
	"fmt"
	"log"

	myPostgres "github.com/NarthurN/FIOapi/internal/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Storage - интерфейс для работы с пользователями в БД.
type Storage interface {
	// RunMigrationsUp() error
	// RunMigrationsDown() error
	// Create(ctx context.Context, user *User) error
	// GetByID(ctx context.Context, id int) (*User, error)
	// GetByEmail(ctx context.Context, email string) (*User, error)
	// Update(ctx context.Context, user *User) error
	// Delete(ctx context.Context, id int) error
}

// PostgresStorage реализует Storage для PostgreSQL.
type PostgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	db, err := myPostgres.New()
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{DB: db}, nil
}

func (ps *PostgresStorage) RunMigrationsUp() error {
	// Создаем экземпляр драйвера для PostgreSQL
	driver, err := postgres.WithInstance(ps.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	// Инициализируем мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("migration initialization failed: %w", err)
	}

	// Применяем миграции up
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not apply migrations: %w", err)
	}

	log.Println("Migrations up applied successfully")
	return nil
}

func (ps *PostgresStorage) RunMigrationsDown() error {
	// Создаем экземпляр драйвера для PostgreSQL
	driver, err := postgres.WithInstance(ps.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	// Инициализируем мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("migration initialization failed: %w", err)
	}

	// Применяем миграции down
	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not apply migrations down: %w", err)
	}

	log.Println("Migrations down applied successfully")
	return nil
}
