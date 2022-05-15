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
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	if err = service.Run(cfg); err != nil {
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
	log.Println("USE DEFAULT CONFIG")
	return &service.Config{
		Auth: service.AuthConfig{Addr: "localhost", Port: 50051},
		DB: service.ConfigDB{
			Username: "postgres",
			Host:     "localhost",
			Port:     5432,
			Dbname:   "store",
			Password: "postgres",
		},
	}
}
