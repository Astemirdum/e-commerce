package interceptor

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"log"
	"strings"
	"time"

	authv1 "github.com/Astemirdum/e-commerce/gen/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthClient struct {
	auth authv1.AuthServiceClient
	log  *zap.Logger
}

func NewAuthClient(addr string, log *zap.Logger) (*AuthClient, error) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &AuthClient{auth: authv1.NewAuthServiceClient(cc), log: log}, nil
}

func (a *AuthClient) AuthInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		a.log.Error("retrieving metadata failed")
		return nil, status.Errorf(codes.Internal, "retrieving metadata failed")
	}
	authMD, ok := md["auth"]
	if !ok {
		a.log.Error("no auth details supplied")
		return nil, status.Errorf(codes.InvalidArgument, "no auth details supplied")
	}
	headerToken := authMD[0]
	if headerToken == "" {
		a.log.Error("empty authMD")
		return nil, status.Errorf(codes.Unauthenticated, "empty authMD")
	}

	headerTokenParts := strings.Split(headerToken, " ")
	if len(headerTokenParts) != 2 || headerTokenParts[0] != "Bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	token := headerTokenParts[1]
	if len(token) == 0 {
		a.log.Error("empty token")
		return nil, status.Errorf(codes.Unauthenticated, "empty token")
	}
	a.log.Info("", zap.String("token", token))
	resp, err := a.auth.Validate(ctx, &authv1.ValidateRequest{Token: token})
	if err != nil {
		a.log.Error("token not valid", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "token not valid %v", err)
	}
	a.log.Info("", zap.Int64("user_id", resp.UserId))
	ctx = context.WithValue(ctx, "user_id", resp.UserId)

	reply, err := handler(ctx, req)

	log.Printf("request - Method:%s  Duration:%s	Error:%v",
		info.FullMethod,
		time.Since(start),
		err)
	return reply, err
}

func GetUserFromCtx(ctx context.Context) (int64, error) {
	id := ctx.Value("user_id")
	userId, ok := id.(int64)
	if !ok {
		return 0, errors.New("no user_id")
	}
	return userId, nil
}
