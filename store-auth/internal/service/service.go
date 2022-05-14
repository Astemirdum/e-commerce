package service

import authv1 "github.com/Astemirdum/e-commerce/gen/auth/v1"

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer
}
