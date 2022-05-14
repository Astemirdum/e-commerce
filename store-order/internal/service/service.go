package service

import orderv1 "github.com/Astemirdum/e-commerce/gen/order/v1"

type OrderServer struct {
	orderv1.UnimplementedOrderServiceServer
}
