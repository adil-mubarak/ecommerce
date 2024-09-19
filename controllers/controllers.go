package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/adil-mubarak/ecommerce/database"
	"github.com/adil-mubarak/ecommerce/models"
	generate "github.com/adil-mubarak/ecommerce/tokens"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DA *gorm.DB = database.GetDB()
var Validate = validator.New()

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := Validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var existingUser models.User
		if err := DA.WithContext(ctx).Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		if err := DA.WithContext(ctx).Where("phone = ?", user.Phone).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is already in use"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		token, refreshToken, _ := generate.TokenGenerator(user.Email, user.FirstName, user.LastName, user.UserID)
		user.Token = token
		user.RefreshToken = refreshToken

		if err := DA.WithContext(ctx).Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not created"})
			return
		}

		c.JSON(http.StatusCreated, "Successfully Signed Up!!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var foundUser models.User
		if err := DA.WithContext(ctx).Where("email = ?", user.Email).First(&foundUser).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login or password incorrect"})
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Login or password incorrect"})
			return
		}

		token, refreshToken, _ := generate.TokenGenerator(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.UserID)
		generate.UpdateAllTokens(token, refreshToken, foundUser.UserID)

		c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": refreshToken})
	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var product models.Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := DA.WithContext(ctx).Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Product not created"})
			return
		}

		c.JSON(http.StatusOK, "Successfully added the product!")
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []models.Product
		if err := DA.Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong. Please try again later."})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProducts []models.Product
		queryParam := c.Query("name")
		if queryParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search query"})
			return
		}

		if err := DA.Where("name LIKE ?", "%"+queryParam+"%").Find(&searchProducts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching the products"})
			return
		}

		c.JSON(http.StatusOK, searchProducts)
	}
}
