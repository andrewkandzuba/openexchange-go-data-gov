package api

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"strings"
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

func TestLoader_NewsToJson_Success(t *testing.T) {
	api, _ := NewCommerceApi(localHost, apiKey)

	body, err := api.News()

	assert.Nil(t, err)
	assert.True(t, strings.Contains(body, "jsonapi"))

	value, dataType, offset, err := jsonparser.Get([]byte(body), "data")

	assert.NotNil(t, value)
	assert.Equal(t, dataType, jsonparser.Array)
	assert.Equal(t, 60423, offset)
	assert.Nil(t, err)

	news := make([]string, 0, 3)
	offset, err = jsonparser.ArrayEach([]byte(body), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		news = append(news, string(value))
	}, "data")

	assert.Nil(t, err)
	assert.Equal(t, 3, len(news))
}

func TestCommerceStream_ConsumeFromSink_Success(t *testing.T) {

	a := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)

	go func() {
		for _, v := range a {
			c <- v
		}
	}()

	cnt := 0
	loop := true

	for loop {
		select {
		case <-c:
			cnt++
		case <-time.After(3 * time.Second):
			fmt.Println("out of time :(")
			loop = false
		}
	}

	assert.Equal(t, len(a), cnt)
}

