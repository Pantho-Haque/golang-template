package stores

import (
	"gorm.io/gorm"
)

type StoreHolder struct {
	UserStore  *UserStore
	// ← add
}

func NewStoreHolder(db *gorm.DB) *StoreHolder {
	return &StoreHolder{
		UserStore: &UserStore{db: db},
		// ← add
	}
}
