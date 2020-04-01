package stream

import (
	"errors"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/connector"
	"github.com/buger/jsonparser"
	"gopkg.in/validator.v2"
	"log"
)

type CommerceStream struct {
	Sink *connector.Connector `validate:"nonnil"`
}

func NewCommerceStream(sink *connector.Connector) (*CommerceStream, error) {
	instance := &CommerceStream{
		Sink: sink,
	}

	if errs := validator.Validate(instance); errs != nil {
		// ToDo: Create a test to handle log.Fatal(...)
		return nil, errors.New(errs.Error())
	}

	return instance, nil
}

func (cs *CommerceStream) Stream() chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		json, err := cs.Sink.News()
		if err != nil {
			// @ToDo: add test to handle this use case.
			log.Fatal(err)
		}

		_, err = jsonparser.ArrayEach([]byte(json), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			log.Println("New article has been found: " + string(value))
			ch <- string(value)
		}, "data")
		if err != nil {
			// @ToDo: add test to handle this use case.
			log.Print(err)
		}

	}()

	return ch
}
