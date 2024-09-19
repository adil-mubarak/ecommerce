package database

import (
	"fmt"
	"log"
	"time"

	"github.com/adil-mubarak/ecommerce/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := "root:kl18jda183079@tcp(localhost:3306)/e_commerse?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Fatal("failed to connect to the database: ",err)
	}
	sqlDB ,err := db.DB()
	if err != nil{
		log.Fatal("failed to get database instance: ",err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 *time.Minute)

	fmt.Println("Successfully connected to the database")

	db.AutoMigrate(&models.User{},&models.Product{},&models.Address{},&models.Order{},models.Payment{},&models.ProductUser{})
	return db
}

func Init(){
	DB = InitDB()
}

func UserData(db *gorm.DB) *gorm.DB{
	return db.Model(&models.User{})
}

func ProductData(db *gorm.DB) *gorm.DB{
	return db.Model(&models.Product{})
}

func GetDB() *gorm.DB {
	return DB
}