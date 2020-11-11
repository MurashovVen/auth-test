package main

import (
	"auth/app"
	"auth/endpoint"
	"auth/model"
	"auth/model/account"
	"auth/model/token"
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Print("loading .env file")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Print("creating datasources")
	account.DataSource = model.GetDataSource("medods", "accounts", context.TODO())
	token.DataSource = model.GetDataSource("medods", "tokens", context.TODO())
	defer account.DataSource.Client.Disconnect(context.TODO())
	defer token.DataSource.Client.Disconnect(context.TODO())

	log.Print("configuring endpoints")
	router := mux.NewRouter()

	router.HandleFunc("/api/account/register", endpoint.Registration).Methods("POST")
	router.HandleFunc("/api/account/login", endpoint.AccessRefreshTokens).Methods("POST")
	router.HandleFunc("/api/account/refresh", endpoint.Refresh).Methods("POST")
	router.HandleFunc("/api/account/refresh", endpoint.DelRefreshToken).Methods("DELETE")
	router.HandleFunc("/api/account/refresh/all", endpoint.DelAllRefreshTokens).Methods("DELETE")

	router.Use(app.JwtAuthentication)

	port := os.Getenv("PORT")

	log.Print("starting server at PORT:" + port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
