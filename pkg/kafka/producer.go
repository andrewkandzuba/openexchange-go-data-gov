package kafka

import (
	"github.com/Shopify/sarama"
	"gopkg.in/validator.v2"
	"strings"
)

type producer struct {
	BootstrapServers string `validate:"nonzero"`
	syncProducer     sarama.SyncProducer
}

func NewKafkaProducer(bootstrapServers string) (*producer, error) {
	instance := &producer{
		BootstrapServers: bootstrapServers,
	}

	if err := validator.Validate(instance); err != nil {
		return nil, err
	}

	if err := instance.open(); err != nil {
		return nil, err
	}

	return instance, nil
}

func (p *producer) open() error {
	var sp, err = sarama.NewSyncProducer(strings.Split(p.BootstrapServers, "\\,"), nil)
	if err != nil {
		return err
	}
	p.syncProducer = sp
	return nil
}

func (p *producer) Send(topic string, data string) error {
	if _, _, err := p.syncProducer.SendMessage(&sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(data)}); err != nil {
		return err
	}
	return nil
}

func (p *producer) Close() error {
	if p.syncProducer != nil {
		if err := p.syncProducer.Close(); err != nil {
			return err
		}
	}
	return nil
}
