package channel

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"gopkg.in/validator.v2"
	"strings"
)

type Producer interface {
	Send(topic string, data string) error
}

type KafkaProducer struct {
	BootstrapServers string `validate:"nonzero"`
}

func NewKafkaProducer(bootstrapServers string) (*KafkaProducer, error) {
	instance := &KafkaProducer{
		BootstrapServers: bootstrapServers,
	}

	if errs := validator.Validate(instance); errs != nil {
		// ToDo: Create a test to handle log.Fatal(...)
		return nil, errors.New(errs.Error())
	}

	return instance, nil
}

func (p *KafkaProducer) Send(topic string, data string) error {
	var sp, err = sarama.NewSyncProducer(strings.Split(p.BootstrapServers, "\\,"), nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := sp.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	_, _, err = sp.SendMessage(&sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(data)})
	if err != nil {
		return err
	}

	return nil
}
