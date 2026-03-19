package stores

import (
	"gorm.io/gorm"
	"pantho/golang/internal/models"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}


func (us *UserStore) GetFirstTenUsers() ([]models.User, error) {
	var users []models.User
	if err := us.db.Limit(10).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserStore) GetUserByPhoneNumber(phoneNumber string, fields []string) (*models.User, error) {
	var user models.User
	result := us.db.Select(fields).Where("number = ?", phoneNumber).Where("type = ?", "user").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (us *UserStore) GetUserById(userID int, fields []string) (*models.User, error) {
	var user models.User
	result := us.db.Select(fields).Where("id = ?", userID).Where("type = ?", "user").First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
