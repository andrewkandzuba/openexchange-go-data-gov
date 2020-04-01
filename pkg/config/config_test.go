package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig_NewConfig_Environment_Success(t *testing.T) {
	_ = os.Setenv("API_ENDPOINT", "my.endpoint.org")

	cfg := NewConfig("testdata/application-test.yaml")

	assert.Equal(t, "my.endpoint.org", cfg.Api.Endpoint)
	assert.Equal(t, "SOME_KEY_FROM_FILE", cfg.Api.Key)
	assert.Equal(t, "sqlite3", cfg.Db.Dialect)
	assert.Equal(t, "test.db", cfg.Db.Host)
	assert.Equal(t, "http://localhost:8080", cfg.Web.Address)
	assert.Equal(t, "localhost:9092", cfg.Kafka.BootstrapServers)

	assert.Equal(t, "test-topic", cfg.Kafka.Consumer.Topics)
	assert.Equal(t, "my-consumer-group", cfg.Kafka.Consumer.Group)
}
