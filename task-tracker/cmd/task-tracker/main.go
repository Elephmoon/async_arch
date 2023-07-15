package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"task-tracker/internal/app"
	"task-tracker/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("cant load env %v", err)
	}

	ctx := context.Background()

	db, err := config.ConnectToPostgres(ctx)
	if err != nil {
		log.Fatalf("cant connect to postgres %v", err)
	}
	defer db.Close(ctx)

	service := app.SetupHTTPService()

	err = http.ListenAndServe(os.Getenv(`HTTP_PORT`), service.Router)
	if err != nil {
		log.Fatalf("get http listen err %v", err)
	}
}
