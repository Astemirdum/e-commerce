package service

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	productv1 "github.com/Astemirdum/e-commerce/gen/product/v1"
	"github.com/Astemirdum/e-commerce/store-product/internal/repo"
	"github.com/Astemirdum/e-commerce/store-product/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Run(cfg *Config) error {
	log := zap.NewExample().Named("product")

	url := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Dbname, cfg.DB.Password)
	db, err := repo.PGConn(url)
	if err != nil {
		log.Error("db conn fail", zap.Error(err))
		return err
	}

	repository := repo.NewProductRepository(db)

	srv := service.NewProductServer(repository, log.Named("service"))

	s := grpc.NewServer()
	addr := fmt.Sprintf("%s:%d", cfg.Product.Addr, cfg.Product.Port)
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	productv1.RegisterProductServiceServer(s, srv)

	log.Info("server has started listen on %s", zap.String("addr", addr))
	go func() {
		if err := s.Serve(ls); err != nil {
			log.Debug("product server stop %v", zap.Error(err))
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Info("graceful shutdown")
	_ = ls.Close()
	s.GracefulStop()

	return nil
}
