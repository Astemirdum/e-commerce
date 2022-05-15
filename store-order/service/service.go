package service

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	orderv1 "github.com/Astemirdum/e-commerce/gen/order/v1"
	"github.com/Astemirdum/e-commerce/store-order/internal/client"
	"github.com/Astemirdum/e-commerce/store-order/internal/repo"
	"github.com/Astemirdum/e-commerce/store-order/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Run(cfg *Config) error {
	log := zap.NewExample().Named("order")

	url := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Dbname, cfg.DB.Password)
	db, err := repo.PGConn(url)
	if err != nil {
		log.Error("db conn fail", zap.Error(err))
		return err
	}

	repository := repo.NewOrderRepository(db)
	pc, err := client.NewProductClientService(fmt.Sprintf("%s:%d", cfg.Product.Addr, cfg.Product.Port))
	if err != nil {
		log.Error("grpc conn product", zap.Error(err))
		return err
	}

	srv := service.NewOrderServer(repository, log.Named("service"), pc)

	s := grpc.NewServer()
	addr := fmt.Sprintf("%s:%d", cfg.Order.Addr, cfg.Order.Port)
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	orderv1.RegisterOrderServiceServer(s, srv)

	log.Info("server has started listen", zap.String("addr", addr))
	go func() {
		if err := s.Serve(ls); err != nil {
			log.Debug("product server stop", zap.Error(err))
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Info("graceful shutdown")
	_ = ls.Close()
	s.GracefulStop()
	_ = pc.Close()

	return nil
}
