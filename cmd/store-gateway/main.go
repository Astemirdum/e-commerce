package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authv1 "github.com/Astemirdum/e-commerce/gen/auth/v1"
	orderv1 "github.com/Astemirdum/e-commerce/gen/order/v1"
	productv1 "github.com/Astemirdum/e-commerce/gen/product/v1"
	"github.com/felixge/httpsnoop"
	"github.com/go-yaml/yaml"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const textArt = "E-commerce GATEWAY"

func main() {
	fmt.Println(textArt)

	cfg := initConfigs()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))

	addr := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s := newServer(cfg)

	log.Printf("server has started listen on %s", addr)
	go func() {
		if err := s.Serve(l); err != nil {
			log.Printf("server stop %v", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("graceful shutdown")
	_ = l.Close()
	ctx, cFn := context.WithTimeout(context.Background(), time.Second)
	defer cFn()
	_ = s.Shutdown(ctx)
}

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
	log.Println("USE DEFAULT CONFIG")
	return &Config{
		ServerConfig{Addr: "localhost", Port: 8080},
	}
}

func newServer(cfg *Config) http.Server {
	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			token := request.Header.Get("Authorization")
			log.Println("token", token)
			return metadata.Pairs("auth", token)
		}),
	)

	if err := authv1.RegisterAuthServiceHandlerFromEndpoint(context.Background(), mux,
		cfg.Auth, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		log.Fatal(err)
	}
	if err := orderv1.RegisterOrderServiceHandlerFromEndpoint(context.Background(), mux,
		cfg.Order, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		log.Fatal(err)
	}
	if err := productv1.RegisterProductServiceHandlerFromEndpoint(context.Background(), mux,
		cfg.Product, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		log.Fatal(err)
	}

	return http.Server{
		Handler: withLogger(mux),
	}
}

func withLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// pass the handler to httpsnoop to get http status and latency
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}
