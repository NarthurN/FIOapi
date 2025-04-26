package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NarthurN/FIOapi/internal/apiclients"
	"github.com/NarthurN/FIOapi/internal/db/migrations"
	"github.com/NarthurN/FIOapi/internal/server"
	"github.com/NarthurN/FIOapi/internal/user"
)

func main() {
	// инициализируем логер
	log := server.SetupLogger("developmentLocal")

	// инициализация базы хранения пользователей
	log.Info("Инициализация базы данных")
	userStorage, err := user.NewStorage()
	if err != nil {
		log.Error("Ошибка инициализации базы данных", "err", err, "op", "main.user.NewStorage()")
		userStorage.DB.Close()
		os.Exit(1)
	}

	defer func() {
		log.Info("Закрытие базы данных...")
		if err := userStorage.DB.Close(); err != nil {
			log.Error("Ошибка закрытия БД", "err", err, "op", "main.userStorage.DB.Close()")
		} else {
			log.Info("Закрыли базу данных.")
		}
	}()

	// запускаем миграции для создания таблиц базы данных
	if err = migrations.RunMigrationsUp(userStorage.DB); err != nil {
		log.Error(
			"Не удалось выполнить миграции up",
			"err", err,
			"op", "main.migrations.RunMigrationsUp(userStorage.DB)",
		)
		os.Exit(1)
	}
	log.Info("Применили миграции.")

	log.Info("Инициализация клиента")
	apiClient := apiclients.New(
		`https://api.agify.io`,
		`https://api.genderize.io`,
		`https://api.nationalize.io`,
		log,
	)

	// запускаем сервис для работы с пользователями
	log.Info("Инициализция сервиса по работе с пользователями userService")
	userService := user.NewService(userStorage, log, apiClient)

	// запускаем сервер
	log.Info("Инициализция веб-сервера")
	userServer := server.Init(userService)

	go func() {
		log.Info("Сервер слушает по адресу http://localhost:8080")
		if err := userServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Ошибка сервера", "err", err, "op", "main.userServer.ListenAndServe()")
			os.Exit(1)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Закрываем сервер...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := userServer.Shutdown(ctx); err != nil {
		log.Error("Ошибка остановки сервера", "err", err, "op", "main.userServer.Shutdown(ctx)")
		os.Exit(1)
	} else {
		log.Info("Сервер успешно остановлен.")
	}

	log.Info("Откатываем миграции...")
	if err = migrations.RunMigrationsDown(userStorage.DB); err != nil {
		log.Error(
			"Не удалось выполнить миграции down",
			"err", err,
			"op", "main.migrations.RunMigrationsDown(userStorage.DB)",
		)
	} else {
		log.Info("Откатили миграции.")
	}
}
