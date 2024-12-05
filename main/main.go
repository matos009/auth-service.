package main

import (
	"auth-service/db"
	"auth-service/handlers"
	"log"
	"net/http"
)

func main() {

	db.Init()

	http.HandleFunc("/issue-tokens", handlers.IssueTokenHandler)
	http.HandleFunc("/refresh-tokens", handlers.RefreshTokensHandler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
