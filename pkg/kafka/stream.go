package kafka

import (
	"errors"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
	"gopkg.in/validator.v2"
)

type Producer struct {
	Topic string `validate:"nonzero"`
	Config *config.Config `validate:"nonzero"`
}

// ToDo: https://github.com/segmentio/kafka-go
func NewProducer(topic string) (*Producer, error) {
	instance := &Producer{
		Topic: topic,
	}

	if errs := validator.Validate(instance) ; errs != nil {
		// ToDo: Create a test to handle log.Fatal(...)
		return nil, errors.New(errs.Error())
	}

	return instance, nil
}

func (p *Producer) produce() error {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{p.Config.Kafka.BootstrapServers},
		Topic:   p.Topic,
		Balancer: kafka.Murmur2Balancer{},
	})
	defer w.Close()

	err := w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{
			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	)
	if err != nil {
		return err
	}

	return nil;
}
