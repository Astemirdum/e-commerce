package service

import productv1 "github.com/Astemirdum/e-commerce/gen/product/v1"

type ProductServer struct {
	productv1.UnimplementedProductServiceServer
}
