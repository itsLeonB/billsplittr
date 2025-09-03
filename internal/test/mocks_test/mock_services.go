package mocks_test

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/dto"
)

type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

func (m *MockAuthService) Register(ctx context.Context, request dto.RegisterRequest) error {
	ret := m.ctrl.Call(m, "Register", ctx, request)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockAuthServiceMockRecorder) Register(ctx, request interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", nil, ctx, request)
}

func (m *MockAuthService) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	ret := m.ctrl.Call(m, "Login", ctx, request)
	ret0, _ := ret[0].(dto.LoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthServiceMockRecorder) Login(ctx, request interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", nil, ctx, request)
}

func (m *MockAuthService) VerifyToken(ctx context.Context, token string) (bool, map[string]any, error) {
	ret := m.ctrl.Call(m, "VerifyToken", ctx, token)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(map[string]any)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (mr *MockAuthServiceMockRecorder) VerifyToken(ctx, token interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", nil, ctx, token)
}

type MockFriendshipService struct {
	ctrl     *gomock.Controller
	recorder *MockFriendshipServiceMockRecorder
}

type MockFriendshipServiceMockRecorder struct {
	mock *MockFriendshipService
}

func NewMockFriendshipService(ctrl *gomock.Controller) *MockFriendshipService {
	mock := &MockFriendshipService{ctrl: ctrl}
	mock.recorder = &MockFriendshipServiceMockRecorder{mock}
	return mock
}

func (m *MockFriendshipService) EXPECT() *MockFriendshipServiceMockRecorder {
	return m.recorder
}

func (m *MockFriendshipService) IsFriends(ctx context.Context, profileID1, profileID2 uuid.UUID) (bool, bool, error) {
	ret := m.ctrl.Call(m, "IsFriends", ctx, profileID1, profileID2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (mr *MockFriendshipServiceMockRecorder) IsFriends(ctx, profileID1, profileID2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFriends", nil, ctx, profileID1, profileID2)
}

func (m *MockFriendshipService) CreateAnonymous(ctx context.Context, request dto.NewAnonymousFriendshipRequest) (dto.FriendshipResponse, error) {
	ret := m.ctrl.Call(m, "CreateAnonymous", ctx, request)
	ret0, _ := ret[0].(dto.FriendshipResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockFriendshipServiceMockRecorder) CreateAnonymous(ctx, request interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAnonymous", nil, ctx, request)
}

func (m *MockFriendshipService) GetAll(ctx context.Context, profileID uuid.UUID) ([]dto.FriendshipResponse, error) {
	ret := m.ctrl.Call(m, "GetAll", ctx, profileID)
	ret0, _ := ret[0].([]dto.FriendshipResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockFriendshipServiceMockRecorder) GetAll(ctx, profileID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", nil, ctx, profileID)
}

type MockProfileService struct {
	ctrl     *gomock.Controller
	recorder *MockProfileServiceMockRecorder
}

type MockProfileServiceMockRecorder struct {
	mock *MockProfileService
}

func NewMockProfileService(ctrl *gomock.Controller) *MockProfileService {
	mock := &MockProfileService{ctrl: ctrl}
	mock.recorder = &MockProfileServiceMockRecorder{mock}
	return mock
}

func (m *MockProfileService) EXPECT() *MockProfileServiceMockRecorder {
	return m.recorder
}

func (m *MockProfileService) GetByID(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error) {
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(dto.ProfileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockProfileServiceMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", nil, ctx, id)
}

func (m *MockProfileService) GetNames(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]string, error) {
	ret := m.ctrl.Call(m, "GetNames", ctx, ids)
	ret0, _ := ret[0].(map[uuid.UUID]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockProfileServiceMockRecorder) GetNames(ctx, ids interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNames", nil, ctx, ids)
}

func (m *MockProfileService) GetByIDs(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]dto.ProfileResponse, error) {
	ret := m.ctrl.Call(m, "GetByIDs", ctx, ids)
	ret0, _ := ret[0].(map[uuid.UUID]dto.ProfileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockProfileServiceMockRecorder) GetByIDs(ctx, ids interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIDs", nil, ctx, ids)
}
