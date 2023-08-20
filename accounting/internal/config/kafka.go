package config

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
)

var brokers = []string{"localhost:9092"}

func NewKafkaClientWithConsumer(ctx context.Context, topic string) (*kgo.Client, error) {
	kafkaClient, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup("accounting"),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		return nil, fmt.Errorf("cant create client %w", err)
	}
	err = kafkaClient.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("cant ping kafka %w", err)
	}

	return kafkaClient, nil
}
