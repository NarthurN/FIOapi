package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NarthurN/FIOapi/internal/db/migrations"
	"github.com/NarthurN/FIOapi/internal/server"
	"github.com/NarthurN/FIOapi/internal/user"
)

func main() {
	// инициализация базы хранения пользователей
	userStorage, err := user.NewStorage()
	if err != nil {
		log.Println(err)
	}

	defer func() {
		log.Println("Закрытие базы данных...")
		if err := userStorage.DB.Close(); err != nil {
			log.Println("Ошибка закрытия БД")
		} else {
			log.Println("Закрыли базу данных.")
		}
	}()

	// запускаем миграции для создания таблиц базы данных
	if err = migrations.RunMigrationsUp(userStorage.DB); err != nil {
		log.Println(err)
	}

	// запускаем сервис для работы с пользователями
	service := user.NewService(userStorage)

	// запускаем сервер
	userServer := server.Init(service)

	go func() {
		log.Println("Server is listening on port 8080...")
		if err := userServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Закрываем сервер...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := userServer.Shutdown(ctx); err != nil {
		log.Fatal("Ошибка остановки сервера:", err)
	} else {
		log.Println("Сервер успешно остановлен.")
	}

	log.Println("Откатываем миграции...")
	if err = migrations.RunMigrationsDown(userStorage.DB); err != nil {
		log.Println(err)
	} else {
		log.Println("Откатили миграции.")
	}
}
