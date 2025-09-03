package mocks_test

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	"github.com/itsLeonB/cocoon-protos/gen/go/auth/v1"
	"google.golang.org/grpc"
)

type MockAuthServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceClientMockRecorder
}

type MockAuthServiceClientMockRecorder struct {
	mock *MockAuthServiceClient
}

func NewMockAuthServiceClient(ctrl *gomock.Controller) *MockAuthServiceClient {
	mock := &MockAuthServiceClient{ctrl: ctrl}
	mock.recorder = &MockAuthServiceClientMockRecorder{mock}
	return mock
}

func (m *MockAuthServiceClient) EXPECT() *MockAuthServiceClientMockRecorder {
	return m.recorder
}

func (m *MockAuthServiceClient) Register(ctx context.Context, in *auth.RegisterRequest, opts ...grpc.CallOption) (*auth.RegisterResponse, error) {
	ret := m.ctrl.Call(m, "Register", ctx, in)
	ret0, _ := ret[0].(*auth.RegisterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthServiceClientMockRecorder) Register(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthServiceClient)(nil).Register), varargs...)
}

func (m *MockAuthServiceClient) Login(ctx context.Context, in *auth.LoginRequest, opts ...grpc.CallOption) (*auth.LoginResponse, error) {
	ret := m.ctrl.Call(m, "Login", ctx, in)
	ret0, _ := ret[0].(*auth.LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthServiceClientMockRecorder) Login(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthServiceClient)(nil).Login), varargs...)
}

func (m *MockAuthServiceClient) VerifyToken(ctx context.Context, in *auth.VerifyTokenRequest, opts ...grpc.CallOption) (*auth.VerifyTokenResponse, error) {
	ret := m.ctrl.Call(m, "VerifyToken", ctx, in)
	ret0, _ := ret[0].(*auth.VerifyTokenResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthServiceClientMockRecorder) VerifyToken(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockAuthServiceClient)(nil).VerifyToken), varargs...)
}
