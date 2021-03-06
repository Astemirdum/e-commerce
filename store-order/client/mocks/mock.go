// Code generated by MockGen. DO NOT EDIT.
// Source: product.go

// Package mock_client is a generated GoMock package.
package mock_client

import (
	context "context"
	reflect "reflect"

	v1 "github.com/Astemirdum/e-commerce/gen/product/v1"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockProductServiceClient is a mock of ProductServiceClient interface.
type MockProductServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockProductServiceClientMockRecorder
}

// MockProductServiceClientMockRecorder is the mock recorder for MockProductServiceClient.
type MockProductServiceClientMockRecorder struct {
	mock *MockProductServiceClient
}

// NewMockProductServiceClient creates a new mock instance.
func NewMockProductServiceClient(ctrl *gomock.Controller) *MockProductServiceClient {
	mock := &MockProductServiceClient{ctrl: ctrl}
	mock.recorder = &MockProductServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductServiceClient) EXPECT() *MockProductServiceClientMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProductServiceClient) Create(ctx context.Context, in *v1.CreateRequest, opts ...grpc.CallOption) (*v1.CreateResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(*v1.CreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockProductServiceClientMockRecorder) Create(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProductServiceClient)(nil).Create), varargs...)
}

// DecreaseStock mocks base method.
func (m *MockProductServiceClient) DecreaseStock(ctx context.Context, in *v1.DecreaseStockRequest, opts ...grpc.CallOption) (*v1.DecreaseStockResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DecreaseStock", varargs...)
	ret0, _ := ret[0].(*v1.DecreaseStockResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DecreaseStock indicates an expected call of DecreaseStock.
func (mr *MockProductServiceClientMockRecorder) DecreaseStock(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecreaseStock", reflect.TypeOf((*MockProductServiceClient)(nil).DecreaseStock), varargs...)
}

// FindOne mocks base method.
func (m *MockProductServiceClient) FindOne(ctx context.Context, in *v1.FindOneRequest, opts ...grpc.CallOption) (*v1.FindOneResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindOne", varargs...)
	ret0, _ := ret[0].(*v1.FindOneResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne.
func (mr *MockProductServiceClientMockRecorder) FindOne(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockProductServiceClient)(nil).FindOne), varargs...)
}
