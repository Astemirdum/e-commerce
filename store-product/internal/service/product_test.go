package service

import (
	"context"
	"errors"
	"testing"

	productv1 "github.com/Astemirdum/e-commerce/gen/product/v1"
	"github.com/Astemirdum/e-commerce/store-product/internal/repo"
	product_mocks "github.com/Astemirdum/e-commerce/store-product/internal/repo/mocks"
	"github.com/Astemirdum/e-commerce/store-product/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProductServer_DecreaseStock(t *testing.T) {
	type input struct {
		req *productv1.DecreaseStockRequest
	}
	type response struct {
		expectedCode codes.Code
		expectedResp *productv1.DecreaseStockResponse
	}
	type mockBehavior func(ctx context.Context,
		p *product_mocks.MockProductRepository, req *productv1.DecreaseStockRequest)

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
			mockBehavior: func(ctx context.Context, r *product_mocks.MockProductRepository, req *productv1.DecreaseStockRequest) {
				p := &models.Product{
					Id:    1,
					Sku:   "lol",
					Stock: 20,
					Price: 100,
				}
				r.EXPECT().FindOne(ctx, req.GetId()).Return(p, nil)
				r.EXPECT().UpdateProduct(ctx, p)
				r.EXPECT().CreateOrderLog(ctx, &repo.StockRequest{
					ProductId: req.GetId(),
					OrderId:   req.GetOrderId(),
					Count:     req.GetCount(),
				}).Return(nil)
			},
			input: input{
				req: &productv1.DecreaseStockRequest{
					Id:      1,
					OrderId: 2,
					Count:   10,
				},
			},
			response: response{
				expectedCode: codes.OK,
				expectedResp: &productv1.DecreaseStockResponse{},
			},
		},
		{
			name: "ko product not found",
			mockBehavior: func(ctx context.Context, r *product_mocks.MockProductRepository, req *productv1.DecreaseStockRequest) {
				p := &models.Product{
					Id:    1,
					Sku:   "lol",
					Stock: 20,
					Price: 100,
				}
				r.EXPECT().FindOne(ctx, req.GetId()).Return(p, errors.New("product not found"))
			},
			input: input{
				req: &productv1.DecreaseStockRequest{
					Id:      1,
					OrderId: 2,
					Count:   10,
				},
			},
			response: response{
				expectedCode: codes.NotFound,
			},
			wantErr: true,
		},
		{
			name: "ko shortage",
			mockBehavior: func(ctx context.Context, r *product_mocks.MockProductRepository, req *productv1.DecreaseStockRequest) {
				p := &models.Product{
					Id:    1,
					Sku:   "lol",
					Stock: 2,
					Price: 100,
				}
				r.EXPECT().FindOne(ctx, req.GetId()).Return(p, nil)
			},
			input: input{
				req: &productv1.DecreaseStockRequest{
					Id:      1,
					OrderId: 2,
					Count:   10,
				},
			},
			response: response{
				expectedCode: codes.InvalidArgument,
			},
			wantErr: true,
		},
		{
			name: "ko orderlog AlreadyExists",
			mockBehavior: func(ctx context.Context, r *product_mocks.MockProductRepository, req *productv1.DecreaseStockRequest) {
				p := &models.Product{
					Id:    1,
					Sku:   "lol",
					Stock: 20,
					Price: 100,
				}
				r.EXPECT().FindOne(ctx, req.GetId()).Return(p, nil)
				r.EXPECT().UpdateProduct(ctx, p)
				r.EXPECT().CreateOrderLog(ctx, &repo.StockRequest{
					ProductId: req.GetId(),
					OrderId:   req.GetOrderId(),
					Count:     req.GetCount(),
				}).Return(repo.ErrAlreadyExists)
			},
			input: input{
				req: &productv1.DecreaseStockRequest{
					Id:      1,
					OrderId: 2,
					Count:   10,
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
			rep := product_mocks.NewMockProductRepository(c)
			tt.mockBehavior(ctx, rep, tt.input.req)

			product := &ProductServer{
				repo: &repo.Repository{rep},
				log:  log,
			}
			resp, err := product.DecreaseStock(ctx, tt.input.req)
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
