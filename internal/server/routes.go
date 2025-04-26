package server

import (
	"net/http"

	"github.com/NarthurN/FIOapi/internal/user"
)

func InitRoutes(userService *user.UserService) *http.ServeMux {
	router := http.NewServeMux()

	router.Handle(`POST /addUser`, userService.AddUser())

	return router
}
