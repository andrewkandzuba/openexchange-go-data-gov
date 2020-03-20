package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/channel"
	"gopkg.in/validator.v2"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type consumer struct {
	BootstrapServers string `validate:"nonzero"`
	Topics string `validate:"nonzero"`
	Group string `validate:"nonzero"`
	listeners sync.Map
	ready chan bool
}

func NewKafkaConsumer(bootstrapServers string, topics string, group string) (*consumer, error) {
	instance := &consumer{
		BootstrapServers: bootstrapServers,
		Topics: topics,
		Group: group,
	}

	if errs := validator.Validate(instance); errs != nil {
		// ToDo: Create a test to handle log.Fatal(...)
		return nil, errors.New(errs.Error())
	}

	return instance, nil
}

func (kc *consumer) Subscribe(topic string, listener *channel.Listener)  {
	kc.listeners.LoadOrStore(topic, listener)
}

func (kc *consumer) start()  {
	var config = sarama.NewConfig()
	config.Version = sarama.MaxVersion

	ready := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(kc.BootstrapServers, ","), kc.Group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}
	defer func() {
		if err = client.Close(); err != nil {
			log.Panicf("Error closing client: %v", err)
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, strings.Split(kc.Topics, ","), kc); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			ready = make(chan bool)
		}
	}()
	<-ready
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
}

func (kc *consumer) Setup(sarama.ConsumerGroupSession) error {
	close(kc.ready)
	return nil
}

func (kc *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (kc *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}
	return nil
}