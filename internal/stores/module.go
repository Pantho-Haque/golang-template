package stores

import "gorm.io/gorm"

type StoreHolder struct {
	UserStore  *UserStore
	CacheStore *CacheStore
}

func NewStoreHolder(db *gorm.DB) *StoreHolder {
	return &StoreHolder{
		UserStore: &UserStore{db: db},
		// cache: &CacheStore{},
	}
}
