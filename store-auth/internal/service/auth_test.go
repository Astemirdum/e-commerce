package service

import (
	"context"
	"errors"
	"testing"

	authv1 "github.com/Astemirdum/e-commerce/gen/auth/v1"
	"github.com/Astemirdum/e-commerce/store-auth/internal/repo"
	repo_mocks "github.com/Astemirdum/e-commerce/store-auth/internal/repo/mocks"
	service_mocks "github.com/Astemirdum/e-commerce/store-auth/internal/service/jwtoken/mocks"
	"github.com/Astemirdum/e-commerce/store-auth/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthServer_Register(t *testing.T) {
	type input struct {
		req     *authv1.RegisterRequest
		userReq *repo.UserRequest
	}
	type response struct {
		expectedCode codes.Code
		expectedResp *authv1.RegisterResponse
	}
	type mockBehavior func(r *repo_mocks.MockAuthRepository, req *repo.UserRequest)

	ctx := context.Background()
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		input        input
		response     response
		wantErr      bool
	}{
		{
			name: "ok",
			mockBehavior: func(r *repo_mocks.MockAuthRepository, req *repo.UserRequest) {
				r.EXPECT().Create(ctx, req).Return(nil)
			},
			input: input{
				req: &authv1.RegisterRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
				userReq: &repo.UserRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
			},
			response: response{
				expectedCode: codes.OK,
				expectedResp: &authv1.RegisterResponse{},
			},
		},
		{
			name: "ko. invalid input",
			mockBehavior: func(r *repo_mocks.MockAuthRepository, req *repo.UserRequest) {
				r.EXPECT().Create(ctx, req).Return(errors.New("error"))
			},
			input: input{
				req: &authv1.RegisterRequest{
					Email:    "123123",
					Password: "123123",
				},
				userReq: &repo.UserRequest{
					Password: "123123",
					Email:    "123123",
				},
			},
			response: response{
				expectedCode: codes.Internal,
			},
			wantErr: true,
		},
		{
			name: "ko. user already exists",
			mockBehavior: func(r *repo_mocks.MockAuthRepository, req *repo.UserRequest) {
				r.EXPECT().Create(ctx, req).Return(repo.ErrAlreadyExists)
			},
			input: input{
				req: &authv1.RegisterRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
				userReq: &repo.UserRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
			},
			response: response{
				expectedCode: codes.AlreadyExists,
			},
			wantErr: true,
		},
	}
	log := zap.NewExample().Named("test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			rep := repo_mocks.NewMockAuthRepository(c)
			tt.mockBehavior(rep, tt.input.userReq)

			auth := &AuthServer{
				repo: &repo.Repository{rep},
				log:  log,
			}
			resp, err := auth.Register(ctx, tt.input.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.response.expectedCode, status.Code(err))
			require.Equal(t, tt.response.expectedResp, resp)
		})
	}
}

func TestAuthServer_Login(t *testing.T) {
	type input struct {
		req     *authv1.LoginRequest
		userReq *repo.UserRequest
	}
	type response struct {
		expectedCode codes.Code
		expectedResp *authv1.LoginResponse
	}
	type mockBehavior func(r *repo_mocks.MockAuthRepository, jwt *service_mocks.MockJwtToken, req *repo.UserRequest)

	ctx := context.Background()
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		input        input
		response     response
		wantErr      bool
	}{
		{
			name: "ok",
			mockBehavior: func(r *repo_mocks.MockAuthRepository, jwt *service_mocks.MockJwtToken, req *repo.UserRequest) {
				var user = models.User{
					Id:       1,
					Email:    "lol@kek",
					Password: repo.GenPasswordHash(req.Password),
				}
				r.EXPECT().Get(ctx, req).Return(&user, nil)
				token := "token"
				jwt.EXPECT().GenerateToken(&user).Return(token, nil)
			},
			input: input{
				req: &authv1.LoginRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
				userReq: &repo.UserRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
			},
			response: response{
				expectedCode: codes.OK,
				expectedResp: &authv1.LoginResponse{Token: "token"},
			},
		},
		{
			name: "ko. wrong password",
			mockBehavior: func(r *repo_mocks.MockAuthRepository, jwt *service_mocks.MockJwtToken, req *repo.UserRequest) {
				var user = models.User{
					Id:       1,
					Email:    "lol@kek",
					Password: repo.GenPasswordHash(req.Password) + "wrong",
				}
				r.EXPECT().Get(ctx, req).Return(&user, nil)
			},
			input: input{
				req: &authv1.LoginRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
				userReq: &repo.UserRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
			},
			response: response{
				expectedCode: codes.Unauthenticated,
			},
			wantErr: true,
		},
		{
			name: "ko. user not found",
			mockBehavior: func(r *repo_mocks.MockAuthRepository, jwt *service_mocks.MockJwtToken, req *repo.UserRequest) {
				r.EXPECT().Get(ctx, req).Return(nil, errors.New("now rows"))
			},
			input: input{
				req: &authv1.LoginRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
				userReq: &repo.UserRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
			},
			response: response{
				expectedCode: codes.NotFound,
			},
			wantErr: true,
		},
		{
			name: "ko. gen token",
			mockBehavior: func(r *repo_mocks.MockAuthRepository, jwt *service_mocks.MockJwtToken, req *repo.UserRequest) {
				var user = models.User{
					Id:       1,
					Email:    "lol@kek",
					Password: repo.GenPasswordHash(req.Password),
				}
				r.EXPECT().Get(ctx, req).Return(&user, nil)
				jwt.EXPECT().GenerateToken(&user).Return("", errors.New("gen token"))
			},
			input: input{
				req: &authv1.LoginRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
				userReq: &repo.UserRequest{
					Email:    "kek",
					Password: "lol@kek",
				},
			},
			response: response{
				expectedCode: codes.Internal,
			},
			wantErr: true,
		},
	}
	log := zap.NewExample().Named("test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			rep := repo_mocks.NewMockAuthRepository(c)
			jwt := service_mocks.NewMockJwtToken(c)
			tt.mockBehavior(rep, jwt, tt.input.userReq)

			auth := &AuthServer{
				repo: &repo.Repository{rep},
				jwt:  jwt,
				log:  log,
			}
			resp, err := auth.Login(ctx, tt.input.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.response.expectedCode, status.Code(err))
			require.Equal(t, tt.response.expectedResp, resp)
		})
	}
}
