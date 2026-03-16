package models

import "time"

type Parcel struct {
	ID                                uint      `gorm:"primaryKey"`
	UserId                            uint      `gorm:"column:user_id"`
	PickupAddressText                 string    `gorm:"column:pickup_address_text"`
	EstimatedPickupAddressLatitude    float64   `gorm:"column:estimated_pickup_address_latitude"`
	EstimatedPickupAddressLongitude   float64   `gorm:"column:estimated_pickup_address_longitude"`
	ReceiverName                      string    `gorm:"column:receiver_name"`
	ReceiverPhoneNumber               string    `gorm:"column:receiver_phone_number"`
	ReceiverAddressText               string    `gorm:"column:receiver_address_text"`
	DetailedAddress                   string    `gorm:"column:detailed_address"`
	EstimatedReceiverAddressLatitude  float64   `gorm:"column:estimated_receiver_address_latitude"`
	EstimatedReceiverAddressLongitude float64   `gorm:"column:estimated_receiver_address_longitude"`
	ParcelType                        uint      `gorm:"column:parcel_type"`
	ParcelComment                     string    `gorm:"column:parcel_comment"`
	ParcelStatus                      string    `gorm:"column:parcel_status"`
	CityId                            uint      `gorm:"column:city_id"`
	CountryId                         uint      `gorm:"column:country_id"`
	Charge                            float64   `gorm:"column:charge"`
	Meta                              string    `gorm:"column:meta"`
	DueMeta                           string    `gorm:"column:due_meta"`
	IsPromoApplied                    uint      `gorm:"column:is_promo_applied"`
	PromoId                           uint      `gorm:"column:promo_id"`
	DriverCut                         float64   `gorm:"column:driver_cut"`
	DriverType                        int32     `gorm:"column:driver_type"`
	Due                               float64   `gorm:"column:due"`
	Distance                          float64   `gorm:"column:distance"`
	Duration                          float64   `gorm:"column:duration"`
	WillPay                           string    `gorm:"column:will_pay"`
	ItemName                          string    `gorm:"column:item_name"`
	ItemPrice                         *float64  `gorm:"column:item_price"`
	IsReviewRequested                 int       `gorm:"column:is_review_requested"`
	OrderId                           string    `gorm:"column:order_id"`
	PaymentType                       uint      `gorm:"column:payment_type"`
	RetryEnabled                      uint      `gorm:"column:retry_enabled"`
	IsPaid                            uint      `gorm:"column:is_paid"`
	CreatedAt                         time.Time `gorm:"column:created_at"`
	UpdatedAt                         time.Time `gorm:"column:updated_at"`
	PickupAddressPlaceId              int64     `gorm:"-"`
	ReceiverAddressPlaceId            int64     `gorm:"-"`
}
