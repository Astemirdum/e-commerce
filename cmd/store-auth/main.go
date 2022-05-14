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

	if err := service.Run(cfg); err != nil {
		log.Fatal(err)
	}
}

func initConfigs() *service.Config {
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
		cfg := &service.Config{}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			log.Fatal(err)
		}
		return cfg
	}
	return &service.Config{
		Auth: service.AuthConfig{Addr: "localhost", Port: 50051},
	}
}
