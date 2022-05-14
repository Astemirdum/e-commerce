package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Astemirdum/e-commerce/store-product/service"
	"github.com/go-yaml/yaml"
)

const textArt = "E-commerce PRODUCT-SERVICE"

func main() {
	fmt.Println(textArt)

	cfg := initConfigs()

	if err := service.Run(cfg.Product); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Product service.ProductConfig `yaml:"product"`
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
		Product: service.ProductConfig{Addr: "localhost", Port: 50052},
	}
}
