package service

import (
	"fmt"
	authv1 "github.com/Astemirdum/e-commerce/gen/auth/v1"
	"github.com/Astemirdum/e-commerce/store-auth/internal/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type AuthConfig struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

func Run(cfg AuthConfig) error {

	s := grpc.NewServer()
	addr := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	authv1.RegisterAuthServiceServer(s, &service.AuthServer{})

	log.Printf("server has started listen on %s", addr)
	go func() {
		if err := s.Serve(ls); err != nil {
			log.Printf("auth server stop %v", err)
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
