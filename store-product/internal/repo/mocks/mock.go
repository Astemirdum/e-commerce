// Code generated by MockGen. DO NOT EDIT.
// Source: repo.go

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	context "context"
	reflect "reflect"

	repo "github.com/Astemirdum/e-commerce/store-product/internal/repo"
	models "github.com/Astemirdum/e-commerce/store-product/models"
	gomock "github.com/golang/mock/gomock"
)

// MockProductRepository is a mock of ProductRepository interface.
type MockProductRepository struct {
	ctrl     *gomock.Controller
	recorder *MockProductRepositoryMockRecorder
}

// MockProductRepositoryMockRecorder is the mock recorder for MockProductRepository.
type MockProductRepositoryMockRecorder struct {
	mock *MockProductRepository
}

// NewMockProductRepository creates a new mock instance.
func NewMockProductRepository(ctrl *gomock.Controller) *MockProductRepository {
	mock := &MockProductRepository{ctrl: ctrl}
	mock.recorder = &MockProductRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductRepository) EXPECT() *MockProductRepositoryMockRecorder {
	return m.recorder
}

// CreateOrderLog mocks base method.
func (m *MockProductRepository) CreateOrderLog(ctx context.Context, req *repo.StockRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrderLog", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrderLog indicates an expected call of CreateOrderLog.
func (mr *MockProductRepositoryMockRecorder) CreateOrderLog(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrderLog", reflect.TypeOf((*MockProductRepository)(nil).CreateOrderLog), ctx, req)
}

// CreateProduct mocks base method.
func (m *MockProductRepository) CreateProduct(ctx context.Context, product *models.Product) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", ctx, product)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockProductRepositoryMockRecorder) CreateProduct(ctx, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockProductRepository)(nil).CreateProduct), ctx, product)
}

// FindOne mocks base method.
func (m *MockProductRepository) FindOne(ctx context.Context, productId int64) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", ctx, productId)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne.
func (mr *MockProductRepositoryMockRecorder) FindOne(ctx, productId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockProductRepository)(nil).FindOne), ctx, productId)
}

// UpdateProduct mocks base method.
func (m *MockProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProduct", ctx, product)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProduct indicates an expected call of UpdateProduct.
func (mr *MockProductRepositoryMockRecorder) UpdateProduct(ctx, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProduct", reflect.TypeOf((*MockProductRepository)(nil).UpdateProduct), ctx, product)
}
