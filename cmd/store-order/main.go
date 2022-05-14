package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Astemirdum/e-commerce/store-order/service"
	"github.com/go-yaml/yaml"
)

const textArt = "E-commerce ORDER-SERVICE"

func main() {
	fmt.Println(textArt)

	cfg := initConfigs()

	if err := service.Run(cfg.Order); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Order service.OrderConfig `yaml:"order"`
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
		Order: service.OrderConfig{Addr: "localhost", Port: 50053},
	}
}
