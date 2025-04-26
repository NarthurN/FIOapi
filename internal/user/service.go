package user

import "net/http"

// "context"
// "errors"

type UserService struct {
	Storage Storage
}

func NewService(storage Storage) *UserService {
	return &UserService{Storage: storage}
}

func (s *UserService) AddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
