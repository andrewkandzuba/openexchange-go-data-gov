package tests

import (
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/channel"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/config"
	"github.com/andrewkandzuba/openexchange-go-data-gov/pkg/kafka"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

func Test_ProduceToKafka_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	_, err := exec.Command("bash", "env/start.sh").Output()
	assert.Nil(t, err)

	defer func() {
		_, _ = exec.Command("bash", "env/stop.sh").Output()
	}()

	cfg := config.NewConfig("testdata/application-ci.yaml")
	cp, err := kafka.NewKafkaProducer(cfg.Kafka.BootstrapServers)
	assert.Nil(t, err)
	assert.NotNil(t, cp)

	err = channel.Producer.Send(cp, "test", "hello!")
	assert.NotNil(t, err)
}
