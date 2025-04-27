package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/NarthurN/FIOapi/internal/apiclients"
	"github.com/NarthurN/FIOapi/internal/interfaces"
)

// Storage - интерфейс для работы с пользователями в БД.
type Storage interface {
	Create(ctx context.Context, user *User) (int, error)
	GetUsers(ctx context.Context, filter *UserFilter, pagination *Pagination) ([]User, error)
	// Update(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int) (int, error)
}

type Enricher interface {
	EnrichUserData(ctx context.Context, name string) (*apiclients.EnrichmentData, error)
}

type UserService struct {
	log     interfaces.Logger
	storage Storage
	enrich  Enricher
}

func NewService(storage Storage, log interfaces.Logger, enrich Enricher) *UserService {
	return &UserService{
		storage: storage,
		log:     log,
		enrich:  enrich,
	}
}

func (s *UserService) AddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var incomingUser UserJSON
		if err := json.NewDecoder(r.Body).Decode(&incomingUser); err != nil {
			s.log.Debug("Невозможно получить данные из тела запроса")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if incomingUser.Name == "" || incomingUser.Surname == "" {
			s.log.Debug("Имя или фамилия не заданы")
			http.Error(w, "Имя или фамилия не заданы", http.StatusBadRequest)
			return
		}

		enrischedData, err := s.enrich.EnrichUserData(r.Context(), "Dmitriy")
		if err != nil {
			s.log.Debug("Нет имени в API")
			http.Error(w, "Введите другое имя", http.StatusBadRequest)
			return
		}

		user := &User{}
		user.Name = incomingUser.Name
		user.Surname = incomingUser.Surname
		user.Patronymic = incomingUser.Patronymic
		user.Age = enrischedData.Age
		user.Sex = enrischedData.Gender
		user.Nationality = enrischedData.Nationality

		id, err := s.storage.Create(context.Background(), user)
		if err != nil {
			s.log.Error("Ошибка при добавлении в базу данных", "err", err, "op", "internal/user/service.go/AddUser()")
			http.Error(w, "Ошибка при добавлении в базу данных", http.StatusInternalServerError)
			return
		}

		user.ID = id
		s.log.Debug(
			"Добавлен User",
			slog.Group("user",
				"id", id,
				"name", user.Name,
				"surname", user.Surname,
				"patronymic", user.Patronymic,
				"age", user.Age,
				"sex", user.Sex,
				"nationality", user.Nationality,
			))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"status": "success",
			"user":   user,
		})
	}
}

func (s *UserService) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := parseFilters(r.URL.Query())
		pagination := parsePagination(r.URL.Query())
		// Получаем пользователей из БД
		users, err := s.storage.GetUsers(r.Context(), filter, pagination)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Формируем ответ
		response := UsersResponse{
			Users:   users,
			Page:    pagination.Page,
			PerPage: pagination.PerPage,
		}

		s.log.Debug("Получены Users", "количество", len(response.Users))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func (s *UserService) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "internal/user/service.go.DeleteUser()"
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.log.Error("Неверный формат id", "err", err, "op", op)
			http.Error(w, "Неверный формат id", http.StatusBadRequest)
			return
		}

		rowsAffected, err := s.storage.DeleteUser(r.Context(), id)
		if err != nil {
			s.log.Error("Ошибка удаления из БД", "err", err, "op", op)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if rowsAffected == 0 {
			s.log.Debug("Пользователь в БД не найден", "op", op)
			http.Error(w, "User не найден", http.StatusNotFound)
			return
		}

		// Формируем успешный ответ
		response := map[string]string{
			"message": fmt.Sprintf("Пользователь с id %d удален", id),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
