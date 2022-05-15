package service

import (
	"context"
	"errors"

	productv1 "github.com/Astemirdum/e-commerce/gen/product/v1"
	"github.com/Astemirdum/e-commerce/store-product/internal/models"
	"github.com/Astemirdum/e-commerce/store-product/internal/repo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ProductServer) Create(ctx context.Context, req *productv1.CreateRequest) (*productv1.CreateResponse, error) {
	productId, err := s.repo.CreateProduct(ctx, &models.Product{
		Price: req.GetPrice(),
		Sku:   req.GetSku(),
		Stock: req.GetStock(),
	})
	if err != nil {
		s.log.Error("create product", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "create product fail: %v", err)
	}
	s.log.Info("product created", zap.Int64("product id", productId))
	return &productv1.CreateResponse{Id: productId}, nil
}

func (s *ProductServer) FindOne(ctx context.Context, req *productv1.FindOneRequest) (*productv1.FindOneResponse, error) {
	product, err := s.repo.FindOne(ctx, req.GetId())
	if err != nil {
		s.log.Error("findOne", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}

	s.log.Info("findOne", zap.Any("product", product))
	return &productv1.FindOneResponse{Product: &productv1.Product{
		Id:    product.Id,
		Sku:   product.Sku,
		Stock: product.Stock,
		Price: product.Price,
	}}, nil
}

func (s *ProductServer) DecreaseStock(ctx context.Context, req *productv1.DecreaseStockRequest) (*productv1.DecreaseStockResponse, error) {
	product, err := s.repo.FindOne(ctx, req.GetId())
	if err != nil {
		s.log.Error("product not found", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "product not found %v", err)
	}

	if product.Stock < req.GetCount() {
		return nil, status.Errorf(codes.InvalidArgument, "stock shortage, remain: %d", product.Stock)
	}
	product.Stock -= req.GetCount()

	if err = s.repo.UpdateProduct(ctx, product); err != nil {
		s.log.Error("update product", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	// Idempotence
	if err = s.repo.CreateOrderLog(ctx, &repo.StockRequest{
		ProductId: req.GetId(),
		OrderId:   req.GetOrderId(),
		Count:     req.GetCount(),
	}); err != nil {
		s.log.Error("createOrderLog", zap.Error(err))
		if errors.Is(err, repo.ErrAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "%v", err)
		}
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	s.log.Info("decreaseStock", zap.Any("orderId", req.GetOrderId()))
	return &productv1.DecreaseStockResponse{}, nil
}
