package main

import (
	"auth/internal/config"
	"auth/internal/pkg/handlers"
	"context"
	"fmt"
	"log"
	"net/http"
)

const (
	port = 8081
)

func main() {
	ctx := context.Background()

	db, err := config.ConnectToPostgres(ctx)
	if err != nil {
		log.Fatalf("cant conect to postgres %v", err)
	}
	defer db.Close(ctx)
	kafkaClient, err := config.NewKafkaClient(ctx)
	if err != nil {
		log.Fatalf("cant connect to kafka %v", err)
	}
	defer kafkaClient.Close()

	httpService := handlers.NewService(db, kafkaClient)

	http.HandleFunc("/login", httpService.Login)
	http.HandleFunc("/signup", httpService.Signup)

	log.Printf("Server is running at %d port.\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
