package service

import (
	authv1 "github.com/Astemirdum/e-commerce/gen/auth/v1"
	"go.uber.org/zap"

	"github.com/Astemirdum/e-commerce/store-auth/internal/repo"
)

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer

	repo *repo.Repository
	jwt  *JwtWrapper
	log  *zap.Logger
}

func NewAuthServer(repo *repo.Repository, jwt *JwtWrapper, log *zap.Logger) *AuthServer {
	return &AuthServer{
		repo: repo,
		jwt:  jwt,
		log:  log,
	}
}
