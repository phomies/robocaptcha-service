package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Main function
func main() {
	godotenv.Load(".env")
	server := gin.New()
	server.POST("/incoming", httpIncoming)
	server.POST("/incoming/:times", httpIncoming)
	server.POST("/verify/:verifyNum/:verifyWord/:times", httpVerify)
	server.GET("/", healthCheck)
	fmt.Println(http.ListenAndServe(":5000", server))
}
