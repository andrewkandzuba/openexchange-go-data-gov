package feed

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

const commerceApiEndpoint = "https://api.commerce.gov/api"

func TestCommerceStream_Ci(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	api, _ := NewCommerceApi(commerceApiEndpoint, apiKey)
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
				if err {
					news = append(news, v)
				}
				break RangeLoop
			case <-time.After(10 * time.Second):
				break RangeLoop
			}
		}
	}()

	wgReceivers.Wait()

	assert.Equal(t, 1, len(news))
}

