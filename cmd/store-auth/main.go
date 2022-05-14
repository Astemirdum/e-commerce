package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Astemirdum/e-commerce/store-auth/service"
	"github.com/go-yaml/yaml"
)

const textArt = "E-commerce AUTH-SERVICE"

func main() {
	fmt.Println(textArt)

	cfg := initConfigs()

	if err := service.Run(cfg.Auth); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Auth service.AuthConfig `yaml:"auth"`
}

func initConfigs() *Config {
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
	return &Config{
		Auth: service.AuthConfig{Addr: "localhost", Port: 50051},
	}
}
