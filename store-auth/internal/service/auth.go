package service

import (
	"context"

	authv1 "github.com/Astemirdum/e-commerce/gen/auth/v1"
	"github.com/Astemirdum/e-commerce/store-auth/internal/repo"
	"go.uber.org/zap"
)

func (s *AuthServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if err := s.repo.Create(ctx, &repo.UserRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword()},
	); err != nil {
		s.log.Error("register user already exists", zap.String("email", req.Email), zap.Error(err))
		return nil, err
	}
	s.log.Info("register user created", zap.String("email", req.Email))
	return &authv1.RegisterResponse{}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	user, err := s.repo.Get(ctx, &repo.UserRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword()},
	)
	if err != nil {
		s.log.Error("login user does not exist", zap.String("email", req.Email), zap.Error(err))
		return nil, err
	}
	token, err := s.jwt.GenerateToken(user)
	if err != nil {
		s.log.Error("generateToken token", zap.String("email", req.Email), zap.Error(err))
		return nil, err
	}
	s.log.Info("login token has issued", zap.String("token", token))
	return &authv1.LoginResponse{Token: token}, nil
}

func (s *AuthServer) Validate(ctx context.Context, req *authv1.ValidateRequest) (*authv1.ValidateResponse, error) {
	claims, err := s.jwt.ParseToken(req.GetToken())
	if err != nil {
		s.log.Error("parseToken invalid", zap.Error(err))
		return nil, err
	}

	user, err := s.repo.Get(ctx, &repo.UserRequest{Email: claims.Email})
	if err != nil {
		s.log.Error("validate user does not exist", zap.String("email", claims.Email), zap.Error(err))
		return nil, err
	}

	s.log.Info("login token has issued", zap.Int64("UserId", user.Id))
	return &authv1.ValidateResponse{UserId: user.Id}, nil
}
