package config

import (
	"github.com/vrischmann/envconfig"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Api struct {
		Endpoint string `yaml:"endpoint" envconfig:"optional"`
		Key string `yaml:"key" envconfig:"optional"`
	} `yaml:"api"`
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
	defer f.Close()

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
	log.Panic(err)
}
