package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/adil-mubarak/ecommerce/models"
	"gorm.io/gorm"
)

var (
	ErrCantFindProduct   = errors.New("can't find product")
	ErrCantDecodeProduct = errors.New("can't decode products")
	ErrUserIDIsNotValid  = errors.New("user ID is not valid")
	ErrCantUpdateUser    = errors.New("cannot add product to cart")
	ErrCantRemoveItem    = errors.New("cannot remove item from cart")
	ErrCantGetItem       = errors.New("cannot get item from cart")
	ErrCantBuyCartItem   = errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context, db *gorm.DB, productID uint, userID uint) error {
	var product models.ProductUser
	if err := db.First(&product, productID).Error; err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}

	// Append product to user's cart
	if err := db.Model(&user).Association("UserCart").Append(&product); err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, db *gorm.DB, productID uint, userID uint) error {
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}

	// Delete product from user's cart
	if err := db.Model(&user).Association("UserCart").Delete(&models.ProductUser{ID: productID}); err != nil {
		log.Println(err)
		return ErrCantRemoveItem
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, db *gorm.DB, userID uint) error {
	var user models.User
	if err := db.Preload("UserCart").First(&user, userID).Error; err != nil {
		log.Println(err)
		return ErrUserIDIsNotValid
	}

	order := models.Order{
		OrderedAt: time.Now(),
		OrderCart: user.UserCart,
		Price:     0,
		PaymentMethod: models.Payment{
			COD: true,
		},
	}

	for _, item := range user.UserCart {
		order.Price += item.Price
	}

	if err := db.Model(&user).Association("Orders").Append(&order); err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}

	if err := db.Model(&user).Association("UserCart").Clear(); err != nil {
		log.Println(err)
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, db *gorm.DB, productID uint, userID uint) error {
	var product models.ProductUser
	if err := db.First(&product, productID).Error; err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	order := models.Order{
		OrderedAt: time.Now(),
		OrderCart: []models.ProductUser{product},
		Price:     product.Price,
		PaymentMethod: models.Payment{
			COD: true,
		},
	}

	if err := db.Model(&models.User{}).Where("id = ?", userID).
		Association("Orders").Append(&order); err != nil {
		log.Println(err)
		return ErrCantUpdateUser
	}
	return nil
}
