package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/channel"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewProducer_Success(t *testing.T) {
	p, err := NewKafkaProducer("")
	assert.Nil(t, p)
	assert.NotNil(t, err)
	assert.Equal(t, "BootstrapServers: zero value", err.Error())

	cfg := config.NewConfig("testdata/application-test.yaml")

	p, err = NewKafkaProducer(cfg.Kafka.BootstrapServers)
	assert.NotNil(t, p)
	assert.Nil(t, err)
}

func Test_Send_Success(t *testing.T) {
	var addr, err = mockNewBroker(t, "test")
	assert.Nil(t, err)
	assert.NotNil(t, addr)

	p, err := NewKafkaProducer(addr)
	assert.Nil(t, err)
	assert.NotNil(t, p)

	err = channel.Producer.Send(p, "test", "hello!")
	assert.Nil(t, err)
}

func mockNewBroker(t *testing.T, topic string) (string, error) {
	seedBroker := sarama.NewMockBroker(t, 1)
	leader := sarama.NewMockBroker(t, 2)

	metadataResponse := new(sarama.MetadataResponse)
	metadataResponse.AddBroker(leader.Addr(), leader.BrokerID())
	metadataResponse.AddTopicPartition(topic, 0, leader.BrokerID(), nil, nil, nil, sarama.ErrNoError)
	seedBroker.Returns(metadataResponse)

	prodSuccess := new(sarama.ProduceResponse)
	prodSuccess.AddTopicPartition(topic, 0, sarama.ErrNoError)
	leader.Returns(prodSuccess)

	return seedBroker.Addr(), nil
}
