package user

import (
    // "context"
    // "errors"
)

type Service struct {
    storage Storage
}

func NewService(storage Storage) *Service {
    return &Service{storage: storage}
}

func (s *Service) Get() string {
	return "privet"
}