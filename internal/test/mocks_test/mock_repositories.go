package mocks_test

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/entity"
)

type MockGroupExpenseRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGroupExpenseRepositoryMockRecorder
}

type MockGroupExpenseRepositoryMockRecorder struct {
	mock *MockGroupExpenseRepository
}

func NewMockGroupExpenseRepository(ctrl *gomock.Controller) *MockGroupExpenseRepository {
	mock := &MockGroupExpenseRepository{ctrl: ctrl}
	mock.recorder = &MockGroupExpenseRepositoryMockRecorder{mock}
	return mock
}

func (m *MockGroupExpenseRepository) EXPECT() *MockGroupExpenseRepositoryMockRecorder {
	return m.recorder
}

func (m *MockGroupExpenseRepository) Create(ctx context.Context, entity *entity.GroupExpense) error {
	ret := m.ctrl.Call(m, "Create", ctx, entity)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockGroupExpenseRepositoryMockRecorder) Create(ctx, entity interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", nil, ctx, entity)
}

func (m *MockGroupExpenseRepository) GetByID(ctx context.Context, id uuid.UUID, preloads ...string) (*entity.GroupExpense, error) {
	ret := m.ctrl.Call(m, "GetByID", ctx, id, preloads)
	ret0, _ := ret[0].(*entity.GroupExpense)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockGroupExpenseRepositoryMockRecorder) GetByID(ctx, id interface{}, preloads ...interface{}) *gomock.Call {
	varargs := append([]interface{}{ctx, id}, preloads...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", nil, varargs...)
}

func (m *MockGroupExpenseRepository) Update(ctx context.Context, entity *entity.GroupExpense) error {
	ret := m.ctrl.Call(m, "Update", ctx, entity)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockGroupExpenseRepositoryMockRecorder) Update(ctx, entity interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", nil, ctx, entity)
}

func (m *MockGroupExpenseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockGroupExpenseRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", nil, ctx, id)
}

func (m *MockGroupExpenseRepository) GetAll(ctx context.Context, filter interface{}, preloads ...string) ([]entity.GroupExpense, error) {
	ret := m.ctrl.Call(m, "GetAll", ctx, filter, preloads)
	ret0, _ := ret[0].([]entity.GroupExpense)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockGroupExpenseRepositoryMockRecorder) GetAll(ctx, filter interface{}, preloads ...interface{}) *gomock.Call {
	varargs := append([]interface{}{ctx, filter}, preloads...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", nil, varargs...)
}

func (m *MockGroupExpenseRepository) SyncParticipants(ctx context.Context, groupExpenseID uuid.UUID, participants []entity.ExpenseParticipant) error {
	ret := m.ctrl.Call(m, "SyncParticipants", ctx, groupExpenseID, participants)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockGroupExpenseRepositoryMockRecorder) SyncParticipants(ctx, groupExpenseID, participants interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncParticipants", nil, ctx, groupExpenseID, participants)
}
