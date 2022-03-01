package main

import (
	"net/http"

	"github.com/joho/godotenv"
)

// Main function
func main() {
	godotenv.Load("../.env")
	http.HandleFunc("/incoming", httpIncoming)
	http.HandleFunc("/verify/", httpVerify)
	http.HandleFunc("/mongo/users", getUsersCollection)
	http.HandleFunc("/", healthcheck)
	http.ListenAndServe(":5000", nil)
}
