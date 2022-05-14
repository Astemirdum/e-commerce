package service

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	productv1 "github.com/Astemirdum/e-commerce/gen/product/v1"
	"github.com/Astemirdum/e-commerce/store-product/internal/service"
	"google.golang.org/grpc"
)

type ProductConfig struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

func Run(cfg ProductConfig) error {

	s := grpc.NewServer()
	addr := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	productv1.RegisterProductServiceServer(s, &service.ProductServer{})

	log.Printf("server has started listen on %s", addr)
	go func() {
		if err := s.Serve(ls); err != nil {
			log.Printf("product server stop %v", err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("graceful shutdown")
	_ = ls.Close()
	s.GracefulStop()

	return nil
}
