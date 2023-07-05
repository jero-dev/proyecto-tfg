package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func hello(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, "Hello World!")
}

func main() {

	router := gin.Default()
	router.GET("/hello", hello)

	err := router.Run("localhost:8080")

	if err != nil {
		log.Fatal("There was an error: " + err.Error())
	}
}
