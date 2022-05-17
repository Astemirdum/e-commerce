package config

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	ServerConfig `yaml:"grpc-gateway"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`

	Auth    string `yaml:"addr-auth"`
	Product string `yaml:"addr-product"`
	Order   string `yaml:"addr-order"`
}

func InitConfigs() *Config {
	var cfgPath string
	flag.StringVar(&cfgPath, "c", "", "config file")
	flag.Parse()

	if cfgPath != "" {
		file, err := os.Open(cfgPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		cfg := &Config{}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			log.Fatal(err)
		}
		return cfg
	}
	log.Println("USE DEFAULT CONFIG")
	return &Config{
		ServerConfig{Addr: "localhost", Port: 8080},
	}
}
