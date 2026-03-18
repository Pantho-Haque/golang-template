package stores

import (
	"gorm.io/gorm"
	"pantho/golang/internal/models"
)

type UserStore interface {
	GetUserByPhoneNumber(phoneNumber string, fields []string) (*models.User, error)
	GetUserById(userID int, fields []string) (*models.User, error)
}
type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{db: db}
}

func (us *userStore) GetUserByPhoneNumber(phoneNumber string, fields []string) (*models.User, error) {
	var user models.User
	result := us.db.Select(fields).Where("number = ?", phoneNumber).Where("type = ?", "user").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (us *userStore) GetUserById(userID int, fields []string) (*models.User, error) {
	var user models.User
	result := us.db.Select(fields).Where("id = ?", userID).Where("type = ?", "user").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
