package server

import (
	"net/http"
	_ "github.com/NarthurN/FIOapi/docs"

	"github.com/NarthurN/FIOapi/internal/interfaces"
	"github.com/NarthurN/FIOapi/internal/middleware"
	"github.com/NarthurN/FIOapi/internal/user"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRoutes(userService *user.UserService, log interfaces.Logger) http.Handler {
	router := http.NewServeMux()

	ml := middleware.New(log)

	router.Handle(`POST /addUser`, userService.AddUser())
	router.Handle(`GET /getUsers`, userService.GetUsers())
	router.Handle(`DELETE /deleteUser/{id}`, userService.DeleteUser())
	router.Handle(`PUT /putUser/{id}`, userService.ChangeUser())
	// Добавляем маршрут для Swagger UI
	router.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // URL к файлу документации
	))

	return ml.Log(router)
}
