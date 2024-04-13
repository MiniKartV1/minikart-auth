package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MiniKartV1/minikart-auth/internal/adapters/app/api"
	"github.com/MiniKartV1/minikart-auth/internal/adapters/core/auth"
	"github.com/MiniKartV1/minikart-auth/internal/adapters/framework/left/rest"
	"github.com/MiniKartV1/minikart-auth/internal/adapters/framework/right/db"
	"github.com/MiniKartV1/minikart-auth/internal/adapters/framework/right/email"
	"github.com/MiniKartV1/minikart-auth/internal/ports"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting main function.")

	defer os.Exit(0)

	var err error
	serverLocation := os.Getenv("SERVER_LOCATION")
	fmt.Println("serverLocation", serverLocation)
	if serverLocation != "http://localhost:3000" {
		fmt.Println("Loading .env file")
		errLoad := godotenv.Load()
		if errLoad != nil {
			log.Fatalf("Error loading .env file: %v", errLoad)
		}
	}
	// portrs
	var coreAdapter ports.AuthenticationPort
	var dbaseAdapter ports.DBPort
	var appAdapter ports.APIPort
	var restAdapter ports.RESTPort
	var emailServiceAdapter ports.EmailServicePort
	DB_URI := os.Getenv("DB_URI")
	fmt.Println("db_uri", DB_URI)
	// dbaseDriver := "mongodb"
	dbaseAdapter = db.NewAdapter(DB_URI)
	fmt.Println("Connected to database")

	if err != nil {
		fmt.Println("Error connecting to database")
	}

	defer dbaseAdapter.CloseDBConnection()

	coreAdapter = auth.NewAdapter()
	emailServiceAdapter = email.NewAdapter()
	appAdapter = api.NewAdapter(coreAdapter, dbaseAdapter, emailServiceAdapter)

	restAdapter = rest.NewAdapter(appAdapter)
	restAdapter.Run()

}
