package server

import (
	"net/http"

	"github.com/NarthurN/FIOapi/internal/interfaces"
	"github.com/NarthurN/FIOapi/internal/user"
)

func Init(userService *user.UserService, log interfaces.Logger) *http.Server {
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: InitRoutes(userService, log),
	}

	return server
}
