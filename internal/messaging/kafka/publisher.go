package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Publisher struct {
	producer *kafka.Producer
}

func NewPublisher(brokers string) (*Publisher, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})

	if err != nil {
		return nil, err
	}

	return &Publisher{producer: p}, nil
}

func (p *Publisher) Publish(
	ctx context.Context,
	topic string,
	key string,
	payload []byte,
) error {
	deliveryChan := make(chan kafka.Event, 1)

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(key),
		Value: payload,
	}, deliveryChan)

	if err != nil {
		return err
	}

	select {
	case e := <-deliveryChan:
		m := e.(*kafka.Message)
		return m.TopicPartition.Error
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *Publisher) Close(timeoutMs int) {
	p.producer.Flush(timeoutMs)
	p.producer.Close()
}
