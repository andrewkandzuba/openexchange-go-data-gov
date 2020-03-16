package kafka

import (
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewKafkaConsumer_Success(t *testing.T) {
	c, err := NewKafkaConsumer("")
	assert.Nil(t, c)
	assert.NotNil(t, err)
	assert.Equal(t, "BootstrapServers: zero value", err.Error())

	cfg := config.NewConfig("testdata/application-test.yaml")
	c, err = NewKafkaConsumer(cfg.Kafka.BootstrapServers)
	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func Test_Subscribes(t *testing.T) {
	var cfg = config.NewConfig("testdata/application-test.yaml")
	var c, err = NewKafkaConsumer(cfg.Kafka.BootstrapServers)
	assert.NotNil(t, c)
	assert.Nil(t, err)
}


