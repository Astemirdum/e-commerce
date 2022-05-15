package client

import (
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
	return &ProductClientService{productv1.NewProductServiceClient(cc)}, nil
}
