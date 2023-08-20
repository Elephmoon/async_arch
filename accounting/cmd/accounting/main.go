package main

import (
	"accounting/internal/config"
	"accounting/internal/pkg/handlers"
	"accounting/internal/pkg/kafka_consumers"
	"accounting/internal/pkg/kafka_handlers"
	"context"
	"fmt"
	"log"
	"net/http"
)

const (
	port = 8083
)

func main() {
	ctx := context.Background()
	db, err := config.ConnectToPostgres(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)
	kafkaClient, err := config.NewKafkaClientWithConsumer(ctx, "user-lifecycle")
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaClient.Close()
	pricingKafkaClient, err := config.NewKafkaClientWithConsumer(ctx, "task-pricing")
	if err != nil {
		log.Fatal(err)
	}
	defer pricingKafkaClient.Close()

	httpService := handlers.NewService(db, kafkaClient)

	http.HandleFunc("/close-day", httpService.CloseDay)

	usersHandler := kafka_handlers.NewUsersKafkaHandler(db)
	taskPricingHandler := kafka_handlers.NewTaskPricingKafkaHandler(db, kafkaClient)
	kafka_consumers.MustCreateUsersKafkaConsumer(ctx, kafkaClient, usersHandler.Handle)
	kafka_consumers.MustCreateUsersKafkaConsumer(ctx, pricingKafkaClient, taskPricingHandler.Handle)

	log.Printf("Server is running at %d port.\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
