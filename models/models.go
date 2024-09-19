package models

import (
	"time"
)

type User struct {
	ID             uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName      string        `json:"first_name" validate:"required,min=2,max=30" gorm:"size:30;not null"`
	LastName       string        `json:"last_name"  validate:"required,min=2,max=30" gorm:"size:30;not null"`
	Password       string        `json:"password"   validate:"required,min=6" gorm:"not null"`
	Email          string        `json:"email"      validate:"email,required" gorm:"unique;not null"`
	Phone          string        `json:"phone"      validate:"required" gorm:"unique;not null"`
	Token          string        `json:"token"`
	RefreshToken   string        `json:"refresh_token"`
	CreatedAt      time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	UserID         string        `json:"user_id" gorm:"not null"`
	UserCart       []ProductUser `json:"usercart" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AddressDetails []Address     `json:"address" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	OrderStatus    []Order       `json:"orders" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Product struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductName string `json:"product_name" gorm:"size:100;not null"`
	Price       uint64 `json:"price" gorm:"not null"`
	Rating      uint8  `json:"rating"`
	Image       string `json:"image"`
}

type ProductUser struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductName string `json:"product_name" gorm:"size:100;not null"`
	Price       int    `json:"price" gorm:"not null"`
	Rating      uint   `json:"rating"`
	Image       string `json:"image"`
	UserID      uint   `json:"user_id"`
}

type Address struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	House   string `json:"house_name" gorm:"size:100;not null"`
	Street  string `json:"street_name" gorm:"size:100;not null"`
	City    string `json:"city_name" gorm:"size:100;not null"`
	Pincode string `json:"pin_code" gorm:"size:10;not null"`
	UserID  uint   `json:"user_id" gorm:"not null"`
}

type Order struct {
	ID            uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderCart     []ProductUser `json:"order_list" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	OrderedAt     time.Time     `json:"ordered_on" gorm:"autoCreateTime"`
	Price         int           `json:"total_price" gorm:"not null"`
	Discount      int           `json:"discount"`
	PaymentMethod Payment       `json:"payment_method" gorm:"embedded"`
	UserID        uint          `json:"user_id" gorm:"not null"`
}

type Payment struct {
	Digital bool `json:"digital" gorm:"not null"`
	COD     bool `json:"cod"     gorm:"not null"`
}

