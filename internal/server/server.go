package server

import (
	"fmt"
	"net/http"

	"github.com/NarthurN/FIOapi/config"
	"github.com/NarthurN/FIOapi/internal/interfaces"
	"github.com/NarthurN/FIOapi/internal/user"
)

func Init(userService *user.UserService, log interfaces.Logger, cfg *config.Config) *http.Server {
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      InitRoutes(userService, log),
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
		IdleTimeout:  cfg.ServerIdleTimeout,
	}

	return server
}
