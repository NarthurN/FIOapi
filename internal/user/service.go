package user

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/NarthurN/FIOapi/internal/apiclients"
	"github.com/NarthurN/FIOapi/internal/interfaces"
)

// Storage - интерфейс для работы с пользователями в БД.
type Storage interface {
	Create(ctx context.Context, user *User) (int, error)
	// GetByID(ctx context.Context, id int) (*User, error)
	// GetByEmail(ctx context.Context, email string) (*User, error)
	// Update(ctx context.Context, user *User) error
	// Delete(ctx context.Context, id int) error
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
		s.log.Debug("Добавлен User", slog.Group("user", "id", id, "name", user.Name))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"status": "success",
			"user":   user,
		})
	}
}
