package stream

import (
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/connector"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func Test_NewInstance_Success(t *testing.T) {
	cfg := config.NewConfig("testdata/application-test.yaml")

	api := &connector.Connector{
		cfg.Api.Endpoint,
		cfg.Api.Key,
	}

	cs, err := NewCommerceStream(api)

	assert.Nil(t, err)
	assert.NotNil(t, cs)
	assert.Equal(t, api, cs.Sink)
}

func Test_NewInstance_SinkIsNil_Failure(t *testing.T) {
	cs, err := NewCommerceStream(nil)

	assert.Nil(t, cs)
	assert.NotNil(t, err)
	assert.Equal(t, "Sink: zero value", err.Error())
}

func Test_ConsumeJsonFromSource_Success(t *testing.T) {
	cfg := config.NewConfig("testdata/application-test.yaml")

	api, _ := connector.NewConnector(cfg.Api.Endpoint, cfg.Api.Key)
	cs, _ := NewCommerceStream(api)

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(1)

	news := make([]string, 0, 3)

	feed := cs.Stream()
	go func() {
		defer wgReceivers.Done()

	RangeLoop:
		for {
			select {
			case v, err := <-feed:
				if !err {
					break RangeLoop
				}
				news = append(news, v)
			case <-time.After(3 * time.Second):
				break RangeLoop
			}
		}
	}()

	wgReceivers.Wait()

	assert.Equal(t, 3, len(news))
}
