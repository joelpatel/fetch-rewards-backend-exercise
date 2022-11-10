package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joelpatel/fetch-rewards-backend-exercise/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/alltransactions", routes.GetAllTransactions()) // extra route for testing purposes (view transactions in order of timestamp)
	router.GET("/all", routes.GetAllPayers())
	router.POST("/add", routes.AddTransaction())
	router.POST("/spend", routes.SpendTransaction())

	router.Run(":" + port)
}
