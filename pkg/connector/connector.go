package connector

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"gopkg.in/validator.v2"
	"time"
)

type Connector struct {
	Endpoint string `validate:"nonzero"`
	ApiKey string `validate:"nonzero"`
}

func NewConnector(endpoint string, apiKey string) (*Connector, error)  {

	instance := &Connector{
		Endpoint: endpoint,
		ApiKey:   apiKey,
	}

	if errs := validator.Validate(instance) ; errs != nil {
		// ToDo: Create a test to handle log.Fatal(...)
		return nil, errors.New(errs.Error())
	}

	return instance, nil
}

func (api*Connector) News() (string, error) {

	client := resty.New()

	// @ToDo: Externalize timeout configuration
	client.SetTimeout(1 * time.Minute)

	// @ToDo: Implement retry and back-off
	resp, err := client.R().
		EnableTrace().
		SetQueryParams(map[string]string{
			"api_key": api.ApiKey,
		}).
		SetHeader("Accept", "application/json").
		Get(api.Endpoint + "/news")

	// ToDo: Create a test to handle log.Fatal(...)
	if err != nil {
		return "", err
	}

	return string(resp.Body()), nil
}