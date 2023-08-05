package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"task-tracker/internal/config"
	"task-tracker/internal/pkg/handlers"
	"task-tracker/internal/pkg/kafka_consumers"
	"task-tracker/internal/pkg/kafka_handlers"
)

const (
	port = 8082
)

func main() {
	ctx := context.Background()
	db, err := config.ConnectToPostgres(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)
	kafkaClient, err := config.NewKafkaClientWithConsumer(ctx, "auth-users-info")
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaClient.Close()

	httpService := handlers.NewService(db, kafkaClient)

	http.HandleFunc("/", httpService.MainPage)
	http.HandleFunc("/create-task", httpService.CreateTask)
	http.HandleFunc("/close-task", httpService.CloseTask)
	http.HandleFunc("/shuffle-tasks", httpService.Shuffle)
	http.HandleFunc("/get-my-tasks", httpService.GetMyTasks)

	usersHandler := kafka_handlers.NewUsersKafkaHandler(db)
	kafka_consumers.MustCreateUsersKafkaConsumer(ctx, kafkaClient, usersHandler.Handle)

	log.Printf("Server is running at %d port.\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
