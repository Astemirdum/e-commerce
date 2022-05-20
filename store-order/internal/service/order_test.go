package service

import (
	"context"
	"testing"

	orderv1 "github.com/Astemirdum/e-commerce/gen/order/v1"
	productv1 "github.com/Astemirdum/e-commerce/gen/product/v1"
	"github.com/Astemirdum/e-commerce/store-order/client"
	product_mocks "github.com/Astemirdum/e-commerce/store-order/client/mocks"
	"github.com/Astemirdum/e-commerce/store-order/internal/repo"
	repo_mocks "github.com/Astemirdum/e-commerce/store-order/internal/repo/mocks"
	"github.com/Astemirdum/e-commerce/store-order/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestOrderServer_CreateOrder(t *testing.T) {
	type input struct {
		ctx context.Context
		req *orderv1.CreateOrderRequest
	}
	type response struct {
		expectedCode codes.Code
		expectedResp *orderv1.CreateOrderResponse
	}
	type mockBehavior func(ctx context.Context, r *repo_mocks.MockOrderRepository,
		p *product_mocks.MockProductServiceClient, req *orderv1.CreateOrderRequest)
	tests := []struct {
		name         string
		mockBehavior mockBehavior
		input        input
		response     response
		wantErr      bool
	}{
		{
			name: "ok",
			mockBehavior: func(ctx context.Context, r *repo_mocks.MockOrderRepository,
				p *product_mocks.MockProductServiceClient, req *orderv1.CreateOrderRequest) {
				ord := &models.Order{
					ProductId: req.GetProductId(),
					UserId:    1,
				}
				r.EXPECT().CreateOrder(ctx, ord).Return(nil)
				p.EXPECT().DecreaseStock(ctx, &productv1.DecreaseStockRequest{
					Id:      req.GetProductId(),
					OrderId: ord.Id,
					Count:   req.GetCount(),
				}).Return(&productv1.DecreaseStockResponse{}, nil)
			},
			input: input{
				ctx: context.WithValue(context.Background(), "user_id", int64(1)),
				req: &orderv1.CreateOrderRequest{
					ProductId: 2,
					Count:     10,
				},
			},
			response: response{
				expectedCode: codes.OK,
				expectedResp: &orderv1.CreateOrderResponse{Id: 0},
			},
		},
		{
			name: "ok. Unauthenticated no user",
			mockBehavior: func(ctx context.Context, r *repo_mocks.MockOrderRepository,
				p *product_mocks.MockProductServiceClient, req *orderv1.CreateOrderRequest) {
			},
			input: input{
				ctx: context.Background(),
				req: &orderv1.CreateOrderRequest{
					ProductId: 2,
					Count:     10,
				},
			},
			response: response{
				expectedCode: codes.Unauthenticated,
			},
			wantErr: true,
		},
		{
			name: "ok. product decreaseStock",
			mockBehavior: func(ctx context.Context, r *repo_mocks.MockOrderRepository,
				p *product_mocks.MockProductServiceClient, req *orderv1.CreateOrderRequest) {
				ord := &models.Order{
					ProductId: req.GetProductId(),
					UserId:    1,
				}
				r.EXPECT().CreateOrder(ctx, ord).Return(nil)
				p.EXPECT().DecreaseStock(ctx, &productv1.DecreaseStockRequest{
					Id:      req.GetProductId(),
					OrderId: ord.Id,
					Count:   req.GetCount(),
				}).Return(&productv1.DecreaseStockResponse{}, status.Errorf(codes.Internal, "lol"))
				ord1 := *ord
				ord1.Failed = true
				r.EXPECT().UpdateOrder(ctx, &ord1).Return(nil)
			},
			input: input{
				ctx: context.WithValue(context.Background(), "user_id", int64(1)),
				req: &orderv1.CreateOrderRequest{
					ProductId: 2,
					Count:     10,
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
			rep := repo_mocks.NewMockOrderRepository(c)
			pr := product_mocks.NewMockProductServiceClient(c)
			tt.mockBehavior(tt.input.ctx, rep, pr, tt.input.req)

			order := &OrderServer{
				repo:          &repo.Repository{rep},
				productClient: &client.ProductClientService{pr},
				log:           log,
			}
			resp, err := order.CreateOrder(tt.input.ctx, tt.input.req)
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
