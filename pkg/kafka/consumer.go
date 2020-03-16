package kafka

import (
	"errors"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/channel"
	"gopkg.in/validator.v2"
	"sync"
)

type consumer struct {
	BootstrapServers string `validate:"nonzero"`
	listeners sync.Map
}

func NewKafkaConsumer(bootstrapServers string) (*consumer, error) {
	instance := &consumer{
		BootstrapServers: bootstrapServers,
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
