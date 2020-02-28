package ci

import (
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	feed2 "github.com/andrewkandzuba/openexchange-go-data-gov/pkg/feed"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestCommerceStream_Ci(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	cfg := config.NewConfig("testdata/application-ci.yaml")

	api, _ := feed2.NewCommerceApi(cfg.Api.Endpoint, cfg.Api.Key)
	cs, _ := feed2.NewCommerceStream(api)

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

