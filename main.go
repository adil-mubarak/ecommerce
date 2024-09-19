package main

import (
	"log"
	"os"

	"github.com/adil-mubarak/ecommerce/controllers"
	"github.com/adil-mubarak/ecommerce/database"
	"github.com/adil-mubarak/ecommerce/middleware"
	"github.com/adil-mubarak/ecommerce/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	database.Init()
	db := database.GetDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(db), database.UserData(db))

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	// router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())
	// router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	log.Fatal(router.Run(":" + port))
}
