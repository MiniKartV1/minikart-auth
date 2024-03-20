package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NareshAtnPLUS/naresh-auth/internal/adapters/app/api"
	"github.com/NareshAtnPLUS/naresh-auth/internal/adapters/core/auth"
	"github.com/NareshAtnPLUS/naresh-auth/internal/adapters/framework/left/rest"
	"github.com/NareshAtnPLUS/naresh-auth/internal/adapters/framework/right/db"
	"github.com/NareshAtnPLUS/naresh-auth/internal/adapters/framework/right/email"
	"github.com/NareshAtnPLUS/naresh-auth/internal/ports"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting main function.")

	defer os.Exit(0)

	var err error
	errLoad := godotenv.Load("../.env")
	if errLoad != nil {
		log.Fatalf("Error loading .env file: %v", err)
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
