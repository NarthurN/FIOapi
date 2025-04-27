package server

import (
	"net/http"

	"github.com/NarthurN/FIOapi/internal/interfaces"
	"github.com/NarthurN/FIOapi/internal/middleware"
	"github.com/NarthurN/FIOapi/internal/user"
)

func InitRoutes(userService *user.UserService, log interfaces.Logger) http.Handler {
	router := http.NewServeMux()

	ml := middleware.New(log)

	router.Handle(`POST /addUser`, userService.AddUser())
	router.Handle(`GET /getUsers`, userService.GetUsers())

	return ml.Log(router)
}
