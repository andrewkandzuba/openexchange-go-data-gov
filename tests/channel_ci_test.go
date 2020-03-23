package tests

import (
	"github.com/Shopify/sarama"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/channel"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/kafka"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func Test_ProduceToKafka_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	cfg := config.NewConfig("testdata/application-ci.yaml")
	var topic = "test"

	kafkaProducer, err := kafka.NewKafkaProducer(cfg.Kafka.BootstrapServers)
	assert.Nil(t, err)
	assert.NotNil(t, kafkaProducer)

	kafkaConsumer, err := kafka.NewKafkaConsumer(cfg.Kafka.BootstrapServers, topic, cfg.Kafka.Consumer.Group)
	assert.Nil(t, err)
	assert.NotNil(t, kafkaConsumer)

	var received = make([]string, 0)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	kafkaConsumer.Subscribe(topic, func(value interface{}) {
		if m, ok := value.(*sarama.ConsumerMessage); ok {
			received = append(received, string(m.Value))
		}
		wg.Done()
	})
	kafkaConsumer.Start()

	err = channel.Producer.Send(kafkaProducer, topic, "hello!")
	assert.Nil(t, err)

	err = channel.Producer.Send(kafkaProducer, topic, "hello again!")
	assert.Nil(t, err)

	err = channel.Producer.Close(kafkaProducer)
	assert.Nil(t, err)

	wg.Wait()

	err = channel.Consumer.Close(kafkaConsumer)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(received))
	assert.Equal(t, "hello!", received[0])
	assert.Equal(t, "hello again!", received[1])
}
