package stores

import (
	"gorm.io/gorm"
	"magic.pathao.com/parcel/prism/internal/models"
)

type UserStore interface {
	GetUserByID(userID int, fields []string) (*models.User, error)
	GetUserByPhoneNumber(phoneNumber string, fields []string) (*models.User, error)
	GetUserById(userID int, fields []string) (*models.User, error)
}
type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{db: db}
}

func (us *userStore) GetUserByID(userID int, fields []string) (*models.User, error) {
	var user models.User
	result := us.db.Select(fields).First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
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
