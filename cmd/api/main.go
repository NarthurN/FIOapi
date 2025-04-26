package main

import (
	"fmt"
	"log"
	"time"

	"github.com/NarthurN/FIOapi/internal/user"
)

func main() {
	storage, err := user.NewPostgresStorage()
	if err != nil {
		log.Println(err)
	}
	defer storage.DB.Close()
	err = storage.RunMigrationsUp()
	if err != nil {
		log.Println(err)
	}
	service := user.NewService(storage)

	fmt.Println(service.Get())

	time.Sleep(10 * time.Second)
	err = storage.RunMigrationsDown()
	if err != nil {
		log.Println(err)
	}
}