package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/adil-mubarak/ecommerce/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Application struct {
	ProductDB *gorm.DB
	UserDB    *gorm.DB
}

func NewApplication(productDB *gorm.DB,userDB *gorm.DB) *Application {
	return &Application{
		ProductDB: productDB,
		UserDB: userDB,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		userQueryID := c.Query("userID")

		productID, err := strconv.Atoi(productQueryID)
		if err != nil || productQueryID == "" {
			log.Println("Product ID is invalid")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Product ID is invalid"})
			return
		}

		userID, err := strconv.Atoi(userQueryID)
		if err != nil || userQueryID == "" {
			log.Println("User ID is invalid")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ID is invalid"})
			return
		}

		// Add product to user cart
		if err := database.AddProductToCart(c.Request.Context(), app.UserDB, uint(productID), uint(userID)); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
			return
		}

		c.JSON(http.StatusOK, "Successfully added to cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("productID")
		userQueryID := c.Query("userID")

		productID, err := strconv.Atoi(productQueryID)
		if err != nil || productQueryID == "" {
			log.Println("Product ID is invalid")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Product ID is invalid"})
			return
		}

		userID, err := strconv.Atoi(userQueryID)
		if err != nil || userQueryID == "" {
			log.Println("User ID is invalid")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ID is invalid"})
			return
		}

		// Remove item from cart
		if err := database.RemoveCartItem(c.Request.Context(), app.UserDB, uint(productID), uint(userID)); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
			return
		}

		c.JSON(http.StatusOK, "Successfully removed from cart")
	}
}

func (app *Application) BuyItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")

		userID, err := strconv.Atoi(userQueryID)
		if err != nil || userQueryID == "" {
			log.Println("User ID is invalid")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ID is invalid"})
			return
		}

		// Buy items from cart
		if err := database.BuyItemFromCart(c.Request.Context(), app.UserDB, uint(userID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place the order"})
			return
		}

		c.JSON(http.StatusOK, "Successfully placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		productQueryID := c.Query("pid")

		userID, err := strconv.Atoi(userQueryID)
		if err != nil || userQueryID == "" {
			log.Println("User ID is invalid")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ID is invalid"})
			return
		}

		productID, err := strconv.Atoi(productQueryID)
		if err != nil || productQueryID == "" {
			log.Println("Product ID is invalid")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Product ID is invalid"})
			return
		}

		// Instant buy product
		if err := database.InstantBuyer(c.Request.Context(), app.UserDB, uint(productID), uint(userID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place the order"})
			return
		}

		c.JSON(http.StatusOK, "Successfully placed the order")
	}
}
