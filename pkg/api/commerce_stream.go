package api

import (
	"errors"
	"gopkg.in/validator.v2"
)

type CommerceStream struct {
	Sink *CommerceApi `validate:"nonnil"`
	Source chan string // @ToDo: implement Apache Kafka Producer as a source
}

func NewCommerceStream(sink *CommerceApi) (*CommerceStream, error) {
	instance := &CommerceStream{
		Sink:   sink,
		Source: make(chan string),
	}

	if errs := validator.Validate(instance) ; errs != nil {
		// ToDo: Create a test to handle log.Fatal(...)
		return nil, errors.New(errs.Error())
	}

	return instance, nil
}