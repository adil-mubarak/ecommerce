package controllers

import (
	"net/http"

	"github.com/adil-mubarak/ecommerce/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		var address models.Address
		if err := c.BindJSON(&address); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := DB.First(&user, "id = ?", userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		var count int64
		DB.Model(&models.Address{}).Where("user_id = ?", userID).Count(&count)

		if count < 2 {
			address.UserID = user.ID
			if err := DB.Create(&address).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add address"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Address added successfully"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add more than 2 addresses"})
		}
	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := DB.Model(&models.Address{}).Where("user_id = ? AND id = ?", userID, 1).
			Updates(models.Address{
				House:   editAddress.House,
				Street:  editAddress.Street,
				City:    editAddress.City,
				Pincode: editAddress.Pincode,
			}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update home address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Home address updated successfully"})
	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid user id"})
			c.Abort()
			return
		}

		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := DB.Model(&models.Address{}).Where("user_id = ? AND id = ?", userID, 2).
			Updates(models.Address{
				House:   editAddress.House,
				Street:  editAddress.Street,
				City:    editAddress.City,
				Pincode: editAddress.Pincode,
			}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update work address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Work address updated successfully"})
	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("id")
		if userID == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid user id"})
			c.Abort()
			return
		}

		if err := DB.Where("user_id = ?", userID).Delete(&models.Address{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
	}
}
