package service

import (
	productv1 "github.com/Astemirdum/e-commerce/gen/product/v1"
	"github.com/Astemirdum/e-commerce/store-product/internal/repo"
	"go.uber.org/zap"
)

type ProductServer struct {
	productv1.UnimplementedProductServiceServer

	repo *repo.Repository
	log  *zap.Logger
}

func NewProductServer(repo *repo.Repository, log *zap.Logger) *ProductServer {
	return &ProductServer{repo: repo, log: log}
}
