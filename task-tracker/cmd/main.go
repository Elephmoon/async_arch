package main

import (
	"context"
	"log"
	"task-tracker/config"

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
}
