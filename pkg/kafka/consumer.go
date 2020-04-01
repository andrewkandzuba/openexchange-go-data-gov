package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"gopkg.in/validator.v2"
	"log"
	"strings"
	"sync"
)

type consumer struct {
	BootstrapServers string `validate:"nonzero"`
	Topics           string `validate:"nonzero"`
	Group            string `validate:"nonzero"`
	listeners        sync.Map
	client           struct {
		ready chan bool
		cg    sarama.ConsumerGroup
		cancel context.CancelFunc
	}
}

func NewKafkaConsumer(bootstrapServers string, topics string, group string) (*consumer, error) {
	instance := &consumer{
		BootstrapServers: bootstrapServers,
		Topics:           topics,
		Group:            group,
	}

	if errs := validator.Validate(instance); errs != nil {
		return nil, errors.New(errs.Error())
	}

	return instance, nil
}

func (kc *consumer) Subscribe(topic string, fn func(value interface{})) {
	kc.listeners.LoadOrStore(topic, fn)
}

func (kc *consumer) Start() {
	var config = sarama.NewConfig()
	config.Version = sarama.MaxVersion

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(kc.BootstrapServers, ","), kc.Group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}
	kc.client.cg = client
	kc.client.ready = make(chan bool)
	kc.client.cancel = cancel

	go func() {
		for {
			if err := client.Consume(ctx, strings.Split(kc.Topics, ","), kc); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			kc.client.ready = make(chan bool)
		}
	}()
	<-kc.client.ready
	log.Println("Sarama consumer up and running!...")
}

func (kc *consumer) Close() error {
	kc.client.cancel()
	if err := kc.client.cg.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
	return nil
}

func (kc *consumer) Setup(sarama.ConsumerGroupSession) error {
	close(kc.client.ready)
	return nil
}

func (kc *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (kc *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		if v, ok := kc.listeners.Load(message.Topic); ok {
			var fn = v.(func(value interface{}))
			fn(message)
		}
		session.MarkMessage(message, "")
	}
	return nil
}
