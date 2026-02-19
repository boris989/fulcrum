package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel"
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
	tracer := otel.Tracer("kafka")

	ctx, span := tracer.Start(ctx, "kafka.Publish")
	defer span.End()

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

func (p *Publisher) Check(ctx context.Context) error {
	_, err := p.producer.GetMetadata(nil, false, 3000)
	return err
}
