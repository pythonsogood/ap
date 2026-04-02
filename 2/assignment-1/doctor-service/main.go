package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pythonsogood/ap-assignment1/doctor/database"
)

func main() {
	router := gin.Default()

	db, err := database.SQLiteConnectDB("doctor-service.db")

	if err != nil {
		panic(err.Error())
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, world!",
		})
	})

	if err := router.Run(":8081"); err != nil {
		panic(err.Error())
	}
}
