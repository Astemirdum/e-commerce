package service

import (
	orderv1 "github.com/Astemirdum/e-commerce/gen/order/v1"
	"github.com/Astemirdum/e-commerce/store-order/internal/client"
	"github.com/Astemirdum/e-commerce/store-order/internal/repo"
	"go.uber.org/zap"
)

type OrderServer struct {
	orderv1.UnimplementedOrderServiceServer

	repo *repo.Repository
	log  *zap.Logger

	productClient *client.ProductClientService
}

func NewOrderServer(repo *repo.Repository, log *zap.Logger, pc *client.ProductClientService) *OrderServer {
	return &OrderServer{repo: repo, log: log, productClient: pc}
}
