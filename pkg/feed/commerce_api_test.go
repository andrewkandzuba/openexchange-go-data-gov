package feed

import (
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	endpoint = "http://127.0.0.1:8080/api"
	apiKey   = "TEST_KEY"
)

func TestCommerceApi_NewInstance_Success(t *testing.T) {

	api, err := NewCommerceApi(endpoint, apiKey)

	assert.NotNil(t, api)
	assert.Nil(t, err)

	assert.Equal(t, endpoint, api.Endpoint)
	assert.Equal(t, apiKey, api.ApiKey)
}

func TestCommerceApi_NewInstance_Failure(t *testing.T) {
	api, err := NewCommerceApi("", "")

	assert.Nil(t, api)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "Endpoint: zero value"))
	assert.True(t, strings.Contains(err.Error(), "ApiKey: zero value"))

	api, err = NewCommerceApi(endpoint, "")

	assert.Nil(t, api)
	assert.NotNil(t, err)
	assert.Equal(t, "ApiKey: zero value", err.Error())

	api, err = NewCommerceApi("", apiKey)

	assert.Nil(t, api)
	assert.NotNil(t, err)
	assert.Equal(t, "Endpoint: zero value", err.Error())
}

func TestCommerceApi_News_Success(t *testing.T) {
	cfg := config.NewConfig("testdata/application-test.yaml")

	api, _ := NewCommerceApi(cfg.Api.Endpoint, cfg.Api.Key)

	body, err := api.News()

	assert.Nil(t, err)
	assert.True(t, strings.Contains(body, "jsonapi"))
}

func TestCommerceApi_NewsWrongHost_Failure(t *testing.T) {
	api, _ := NewCommerceApi(endpoint, apiKey)

	body, err := api.News()

	assert.NotNil(t, err)
	assert.Equal(t, 0, len(body))
}
