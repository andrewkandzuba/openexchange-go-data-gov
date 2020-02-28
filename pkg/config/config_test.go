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
}
