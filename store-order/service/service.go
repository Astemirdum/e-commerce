package service

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	orderv1 "github.com/Astemirdum/e-commerce/gen/order/v1"
	"github.com/Astemirdum/e-commerce/store-order/internal/service"
	"google.golang.org/grpc"
)

func Run(cfg *Config) error {

	s := grpc.NewServer()
	addr := fmt.Sprintf("%s:%d", cfg.Order.Addr, cfg.Order.Port)
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	orderv1.RegisterOrderServiceServer(s, &service.OrderServer{})

	log.Printf("server has started listen on %s", addr)
	go func() {
		if err := s.Serve(ls); err != nil {
			log.Printf("order server stop %v", err)
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
