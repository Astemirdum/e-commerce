package service

import (
	"context"

	orderv1 "github.com/Astemirdum/e-commerce/gen/order/v1"
	productv1 "github.com/Astemirdum/e-commerce/gen/product/v1"
	"github.com/Astemirdum/e-commerce/store-order/internal/models"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	order := &models.Order{
		ProductId: req.GetProductId(),
		UserId:    req.GetUserId(),
	}
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		s.log.Error("create order", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "create product fail: %v", err)
	}
	s.log.Info("order created", zap.Any("order", order))

	if _, err := s.productClient.DecreaseStock(ctx, &productv1.DecreaseStockRequest{
		Id:      req.GetProductId(),
		OrderId: order.Id,
		Count:   req.GetCount(),
	}); err != nil {
		s.log.Error("decreaseStock", zap.Error(err))
		order.Failed = true
		if err := s.repo.UpdateProduct(ctx, order); err != nil {
			s.log.Error("updateProduct", zap.Error(err))
		}
		return nil, status.Errorf(status.Code(err), "decreaseStock fail: %v", err)
	}
	s.log.Info("decreaseStock", zap.Any("order", order))
	return &orderv1.CreateOrderResponse{Id: order.Id}, nil
}
