package models

import (
	"time"
)

type User struct {
	Id                       uint      `gorm:"column:id"`
	Uuid                     string    `gorm:"column:uuid"`
	Name                     string    `gorm:"column:name"`
	Email                    string    `gorm:"column:email"`
	Number                   string    `gorm:"column:number"`
	ProfilePic               string    `gorm:"column:profile_pic"`
	Type                     string    `gorm:"column:type"`
	AltNumber                string    `gorm:"column:alt_number"`
	Address                  string    `gorm:"column:address"`
	DateOfBirth              string    `gorm:"column:date_of_birth"`
	Gender                   uint      `gorm:"column:gender"`
	NeedsToReadTos           bool      `gorm:"column:needs_to_read_tos"`
	FirstKnownLatLong        string    `gorm:"column:first_known_latlng"`
	LastKnownLatLong         string    `gorm:"column:last_known_latlng"`
	City                     string    `gorm:"column:city"`
	AlternateEmail           string    `gorm:"column:alternate_email"`
	IsEmailVerified          bool      `gorm:"column:is_email_verified"`
	FirstKnownCityId         uint      `gorm:"column:first_known_city_id"`
	LastKnownCityId          uint      `gorm:"column:last_known_city_id"`
	FirstKnownCountryId      uint      `gorm:"column:first_known_country_id"`
	LastKnownCountryId       uint      `gorm:"column:last_known_country_id"`
	CountryId                uint      `gorm:"column:country_id"`
	CityId                   uint      `gorm:"column:city_id"`
	DriverType               uint      `gorm:"column:driver_type"`
	IsSuspended              bool      `gorm:"column:is_suspended"`
	WalletId                 string    `gorm:"column:wallet_id"`
	WalletIsActive           bool      `gorm:"column:wallet_is_active"`
	Created                  uint      `gorm:"column:created"`
	Updated                  uint      `gorm:"column:updated"`
	Deleted                  uint      `gorm:"column:deleted"`
	CreatedAt                time.Time `gorm:"column:created_at"`
	UpdatedAt                time.Time `gorm:"column:updated_at"`
	DeletedAt                time.Time `gorm:"column:deleted_at"`
	IsOpenForParcels         bool      `gorm:"column:is_open_for_parcels"`
	ParcelRating             string    `gorm:"column:parcel_rating"`
	ParcelRatingParticipants uint      `gorm:"column:parcel_rating_participants"`
}

func (user *User) GetCityID() uint {
	if user.LastKnownCityId != 0 {
		return user.LastKnownCityId
	}
	return 1
}

func (user *User) GetCountryID() uint {
	if user.LastKnownCountryId != 0 {
		return user.LastKnownCountryId
	}
	return 1
}
