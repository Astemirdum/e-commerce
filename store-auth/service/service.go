package service

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	authv1 "github.com/Astemirdum/e-commerce/gen/auth/v1"
	"github.com/Astemirdum/e-commerce/store-auth/internal/repo"
	"github.com/Astemirdum/e-commerce/store-auth/internal/service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Run(cfg *Config) error {
	log := zap.NewExample().Named("auth")

	url := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Dbname, cfg.DB.Password)
	db, err := repo.PGConn(url)
	if err != nil {
		log.Error("db conn fail", zap.Error(err))
		return err
	}

	repository := repo.NewRepository(db)

	srv := service.NewAuthServer(
		repository,
		service.NewJwtWrapper(cfg.Auth.SecretKey, "store-auth", 15),
		zap.NewExample().Named("service"),
	)

	s := grpc.NewServer()
	addr := fmt.Sprintf("%s:%d", cfg.Auth.Addr, cfg.Auth.Port)
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	authv1.RegisterAuthServiceServer(s, srv)

	log.Info("server has started listen on %s", zap.String("addr", addr))
	go func() {
		if err := s.Serve(ls); err != nil {
			log.Debug("auth server stop %v", zap.Error(err))
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
