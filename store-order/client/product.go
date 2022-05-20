package client

import (
	"context"
	productv1 "github.com/Astemirdum/e-commerce/gen/product/v1"
	"google.golang.org/grpc"
)

type ProductClientService struct {
	productv1.ProductServiceClient
}

func NewProductClientService(addr string) (*ProductClientService, error) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &ProductClientService{
		ProductServiceClient: productv1.NewProductServiceClient(cc),
	}, nil
}

//go:generate mockgen -source=product.go -destination=mocks/mock.go

type ProductServiceClient interface {
	Create(ctx context.Context, in *productv1.CreateRequest, opts ...grpc.CallOption) (*productv1.CreateResponse, error)
	FindOne(ctx context.Context, in *productv1.FindOneRequest, opts ...grpc.CallOption) (*productv1.FindOneResponse, error)
	DecreaseStock(ctx context.Context, in *productv1.DecreaseStockRequest, opts ...grpc.CallOption) (*productv1.DecreaseStockResponse, error)
}
