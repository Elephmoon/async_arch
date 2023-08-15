package kafka_consumers

import (
	"context"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
)

type KafkaHandler func(ctx context.Context, kafkaMsgValue []byte) error

func MustCreateUsersKafkaConsumer(ctx context.Context, kgo *kgo.Client, handler KafkaHandler) {
	go consumption(ctx, kgo, handler)
}

func consumption(ctx context.Context, kgo *kgo.Client, handler KafkaHandler) {
	for {
		fetches := kgo.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Printf("get error during kafka consumption %v", errs)
			continue
		}
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			err := handler(ctx, record.Value)
			if err != nil {
				continue
			}
		}
	}
}
