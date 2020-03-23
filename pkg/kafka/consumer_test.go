package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func Test_NewKafkaConsumer_Success(t *testing.T) {
	c, err := NewKafkaConsumer("", "", "")
	assert.Nil(t, c)
	assert.NotNil(t, err)

	c, err = NewKafkaConsumer("localhost:9092", "", "")
	assert.Nil(t, c)
	assert.NotNil(t, err)

	c, err = NewKafkaConsumer("localhost:9092", "test1,test2", "")
	assert.Nil(t, c)
	assert.NotNil(t, err)

	cfg := config.NewConfig("testdata/application-test.yaml")
	c, err = NewKafkaConsumer(cfg.Kafka.BootstrapServers, cfg.Kafka.Consumer.Topics, cfg.Kafka.Consumer.Group)
	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func Test_Subscribes(t *testing.T) {
	var cfg = config.NewConfig("testdata/application-test.yaml")
	var c, err = NewKafkaConsumer(cfg.Kafka.BootstrapServers, cfg.Kafka.Consumer.Topics, cfg.Kafka.Consumer.Group)
	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func Test_Consumer(t *testing.T) {
	consumer := mocks.NewConsumer(t, nil)
	defer func() {
		if err := consumer.Close(); err != nil {
			t.Error(err)
		}
	}()

	var topic = "test"
	consumer.SetTopicMetadata(map[string][]int32{
		topic: {0, 1},
	})
	consumer.ExpectConsumePartition(topic, 0, sarama.OffsetOldest).YieldMessage(&sarama.ConsumerMessage{Value: []byte("hello world")})
	consumer.ExpectConsumePartition(topic, 1, sarama.OffsetOldest).YieldMessage(&sarama.ConsumerMessage{Value: []byte("hello world again")})

	partitionList, _ := consumer.Partitions(topic)
	messages := make([]*sarama.ConsumerMessage, 0)
	initialOffset := sarama.OffsetOldest

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(2)

	for _, partition := range partitionList {
		pc, _ := consumer.ConsumePartition(topic, partition, initialOffset)
		go func(pc sarama.PartitionConsumer) {
			defer func() {
				wgReceivers.Done()
			}()
			feed := pc.Messages()
		RangeLoop:
			for {
				select {
				case v, err := <-feed:
					if !err {
						break RangeLoop
					}
					messages = append(messages, v)
				case <-time.After(3 * time.Second):
					break RangeLoop
				}
			}
		}(pc)
	}

	wgReceivers.Wait()

	assert.True(t, len(messages) > 0)
}
