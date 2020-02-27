package feed

import (
	"bytes"
	"encoding/json"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/model"
	"github.com/stretchr/testify/assert"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCommerceStream_NewInstance_Success(t *testing.T) {
	api := &CommerceApi{
		endpoint,
		apiKey,
	}

	cs, err := NewCommerceStream(api)

	assert.Nil(t, err)
	assert.NotNil(t, cs)
	assert.Equal(t, api, cs.Sink)
}

func TestNewCommerceStream_NewInstance_SinkIsNil_Failure(t *testing.T) {
	cs, err := NewCommerceStream(nil)

	assert.Nil(t, cs)
	assert.NotNil(t, err)
	assert.Equal(t, "Sink: zero value", err.Error())
}

func TestCommerceStream_ConsumeJsonFromSource_Success(t *testing.T) {
	api, _ := NewCommerceApi(localHost, apiKey)
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

	var a model.Article
	r := bytes.NewReader([]byte(news[0]))
	err := json.NewDecoder(r).Decode(&a)

	assert.Nil(t, err)
	assert.Equal(t, "news", a.Type)
	assert.Equal(t, 48, a.Nid)
	assert.Equal(t, "Former Secretary of Commerce Pritzker’s Official Portrait Unveiled", a.Label)
	assert.Equal(t, 1508360246, a.Created)
	assert.Equal(t, 1513969347, a.Update)
	assert.Equal(t, "https://www.commerce.gov/news/press-releases/2017/10/former-secretary-commerce-pritzkers-official-portrait-unveiled", a.Href)
	assert.True(t, strings.Contains(a.Body, "Commerce Wilbur Ross attended the unveiling of former"))
	assert.Equal(t, "FOR IMMEDIATE RELEASE", a.Status)
	assert.Equal(t, "86e95078-92cf-4ff7-a55d-cb2fbffd8b85", a.UUID)
	assert.NotNil(t, a.AdminOfficials)
	assert.Equal(t, 1, len(a.AdminOfficials))
	assert.Equal(t, "9", a.AdminOfficials[0].Id)
	assert.Equal(t, "Wilbur Ross", a.AdminOfficials[0].Label)
	assert.Equal(t, "https://www.commerce.gov/about/leadership/wilbur-ross", a.AdminOfficials[0].Href)
}
