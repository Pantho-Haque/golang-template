package stores

import (
	"gorm.io/gorm"
	"magic.pathao.com/parcel/prism/internal/models"
)

type ParcelStore interface {
	UpdateParcel(parcel *models.Parcel) (*models.Parcel, error)
}
type parcelStore struct {
	db *gorm.DB
}

func NewParcelStore(db *gorm.DB) ParcelStore {
	return &parcelStore{db: db}
}

func (ps *parcelStore) UpdateParcel(parcel *models.Parcel) (*models.Parcel, error) {
	// update the parcel with the new data
	result := ps.db.Model(parcel).Updates(parcel)
	if result.Error != nil {
		return nil, result.Error
	}

	return parcel, nil
}
