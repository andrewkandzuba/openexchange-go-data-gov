package config

import (
	"github.com/vrischmann/envconfig"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Api struct {
		Endpoint string `yaml:"endpoint" envconfig:"optional"`
		Key string `yaml:"key" envconfig:"optional"`
	} `yaml:"api"`
	Db struct{
		Dialect string `yaml:"dialect" envconfig:"optional"`
		Host string `yaml:"host" envconfig:"optional"`
	} `yaml:"db"`
	Web struct{
		Address string `yaml:"address" envconfig:"optional"`
	} `yaml:"web"`
	Kafka struct{
		BootstrapServers string `yaml:"bootstrap-servers" envconfig:"optional"`
		Consumer struct{
			Topics string `yaml:"topics" envconfig:"optional"`
			Group string `yaml:"group" envconfig:"optional"`
		} `yaml:"consumer"`
	} `yaml:"kafka"`
}

func NewConfig(file string) *Config {
	var cfg Config
	readFile(&cfg, file)
	readEnv(&cfg)
	return &cfg
}

func readFile(cfg *Config, file string) {
	f, err := os.Open(file)
	if err != nil {
		processError(err)
	}
	defer func() {
		_ = f.Close()
	}()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Init(cfg)
	if err != nil {
		processError(err)
	}
}

func processError(err error) {
	panic(err)
}
