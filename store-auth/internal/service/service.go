package service

import (
	authv1 "github.com/Astemirdum/e-commerce/gen/auth/v1"
	"github.com/Astemirdum/e-commerce/store-auth/internal/repo"
	"github.com/Astemirdum/e-commerce/store-auth/internal/service/jwtoken"
	"go.uber.org/zap"
)

type AuthServer struct {
	authv1.UnimplementedAuthServiceServer

	repo *repo.Repository
	jwt  jwtoken.JwtToken
	log  *zap.Logger
}

func NewAuthServer(repo *repo.Repository, jwt *jwtoken.JwtWrapper, log *zap.Logger) *AuthServer {
	return &AuthServer{
		repo: repo,
		jwt:  jwt,
		log:  log,
	}
}
