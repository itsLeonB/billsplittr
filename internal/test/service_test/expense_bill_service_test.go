package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/billsplittr/internal/appconstant"
	"github.com/itsLeonB/billsplittr/internal/dto"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/service"
	"github.com/itsLeonB/billsplittr/internal/test/mocks"
	"github.com/itsLeonB/go-crud"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// Mock logger for testing
type mockLogger struct{}

func (m *mockLogger) Debug(args ...interface{})                 {}
func (m *mockLogger) Debugf(format string, args ...interface{}) {}
func (m *mockLogger) Info(args ...interface{})                  {}
func (m *mockLogger) Infof(format string, args ...interface{})  {}
func (m *mockLogger) Warn(args ...interface{})                  {}
func (m *mockLogger) Warnf(format string, args ...interface{})  {}
func (m *mockLogger) Error(args ...interface{})                 {}
func (m *mockLogger) Errorf(format string, args ...interface{}) {}
func (m *mockLogger) Fatal(args ...interface{})                 {}
func (m *mockLogger) Fatalf(format string, args ...interface{}) {}

func TestExpenseBillService_Upload_ValidationErrorInvalidFileType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExpenseBillRepo := mocks.NewMockExpenseBillRepository(ctrl)
	mockStorageRepo := mocks.NewMockStorageRepository(ctrl)
	logger := &mockLogger{}

	svc := service.NewExpenseBillService(mockExpenseBillRepo, mockStorageRepo, "test-bucket", logger)

	request := &dto.UploadBillRequest{
		PayerProfileID:   uuid.New(),
		CreatorProfileID: uuid.New(),
		ImageData:        []byte("test data"),
		ContentType:      "application/pdf", // Invalid content type
		Filename:         "test.pdf",
		FileSize:         1024,
	}

	_, err := svc.Upload(context.Background(), request)

	assert.Error(t, err)
}

func TestExpenseBillService_Upload_ValidationErrorFileTooLarge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExpenseBillRepo := mocks.NewMockExpenseBillRepository(ctrl)
	mockStorageRepo := mocks.NewMockStorageRepository(ctrl)
	logger := &mockLogger{}

	svc := service.NewExpenseBillService(mockExpenseBillRepo, mockStorageRepo, "test-bucket", logger)

	request := &dto.UploadBillRequest{
		PayerProfileID:   uuid.New(),
		CreatorProfileID: uuid.New(),
		ImageData:        make([]byte, appconstant.MaxFileSize+1), // Too large
		ContentType:      "image/jpeg",
		Filename:         "test.jpg",
		FileSize:         appconstant.MaxFileSize + 1,
	}

	_, err := svc.Upload(context.Background(), request)

	assert.Error(t, err)
}

func TestExpenseBillService_GetSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExpenseBillRepo := mocks.NewMockExpenseBillRepository(ctrl)
	mockStorageRepo := mocks.NewMockStorageRepository(ctrl)
	logger := &mockLogger{}

	svc := service.NewExpenseBillService(mockExpenseBillRepo, mockStorageRepo, "test-bucket", logger)

	billID := uuid.New()
	profileID := uuid.New()

	expectedBill := entity.ExpenseBill{
		BaseEntity:       crud.BaseEntity{ID: billID},
		PayerProfileID:   profileID,
		CreatorProfileID: profileID,
		ImageName:        "test-image.jpg",
	}

	mockExpenseBillRepo.EXPECT().
		FindFirst(gomock.Any(), gomock.Any()).
		Return(expectedBill, nil)

	result, err := svc.Get(context.Background(), billID, profileID)

	assert.NoError(t, err)
	assert.Equal(t, billID, result.ID)
	assert.Equal(t, profileID, result.PayerProfileID)
	assert.Equal(t, "test-image.jpg", result.Filename)
}

func TestExpenseBillService_GetNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExpenseBillRepo := mocks.NewMockExpenseBillRepository(ctrl)
	mockStorageRepo := mocks.NewMockStorageRepository(ctrl)
	logger := &mockLogger{}

	svc := service.NewExpenseBillService(mockExpenseBillRepo, mockStorageRepo, "test-bucket", logger)

	billID := uuid.New()
	profileID := uuid.New()

	// Return zero value entity (not found)
	mockExpenseBillRepo.EXPECT().
		FindFirst(gomock.Any(), gomock.Any()).
		Return(entity.ExpenseBill{}, nil)

	_, err := svc.Get(context.Background(), billID, profileID)

	assert.Error(t, err)
}

func TestExpenseBillService_GetUnauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExpenseBillRepo := mocks.NewMockExpenseBillRepository(ctrl)
	mockStorageRepo := mocks.NewMockStorageRepository(ctrl)
	logger := &mockLogger{}

	svc := service.NewExpenseBillService(mockExpenseBillRepo, mockStorageRepo, "test-bucket", logger)

	billID := uuid.New()
	profileID := uuid.New()

	// Return zero value entity (not found due to different creator profile ID)
	mockExpenseBillRepo.EXPECT().
		FindFirst(gomock.Any(), gomock.Any()).
		Return(entity.ExpenseBill{}, nil)

	_, err := svc.Get(context.Background(), billID, profileID)

	assert.Error(t, err)
}
