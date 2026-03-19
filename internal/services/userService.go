package services

import (
	"pantho/golang/internal/models"
	"pantho/golang/internal/stores"

	"go.uber.org/zap"
)

type UserService struct {
	userStore *stores.UserStore
	log *zap.Logger
}

func NewUserService(userStore *stores.UserStore,log *zap.Logger) *UserService {
	return &UserService{
		userStore: userStore,
		log: log,
	}
}

func (s *UserService) GetUsers() (*[]models.User, error) {
	users , err := s.userStore.GetFirstTenUsers();
	if  err != nil {
		return nil, err
	}
	return &users, nil
}
