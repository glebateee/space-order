package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string      `yaml:"env"          env-default:"local"`
	StoragePath string      `yaml:"storage_path" env-required:"true"`
	HttpConfig  *HttpConfig `yaml:"http_config"`
}

type HttpConfig struct {
	Host        string        `yaml:"host"          env-default:"localhost"`
	Port        int           `yaml:"port"          env-default:"2282"`
	Timeout     time.Duration `yaml:"timeout"       env-default:"1h"`
	IdleTimeout time.Duration `yaml:"idle_timeout"  env-default:"2h"`
}

var emptyPath = ""

func MustLoad() *Config {
	configPath := fetchConigPath()
	if configPath == emptyPath {
		panic("config path not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not found" + configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("error reading config file" + err.Error())
	}

	return &cfg
}

func fetchConigPath() string {
	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == emptyPath {
		configPath = os.Getenv("CONFIG_PATH")
	}
	return configPath
}
