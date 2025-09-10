package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/billsplittr/internal/repository"
	"github.com/stretchr/testify/assert"
)

const (
	testKey    = "test-key"
	testBucket = "test-bucket"
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

func TestStorageRepositoryInterface(t *testing.T) {
	// Test that the interface is properly defined
	var repo repository.StorageRepository
	
	// Test interface methods exist (nil interface is fine for this test)
	assert.Nil(t, repo)
}

func TestStorageUploadRequestValidation(t *testing.T) {
	req := &entity.StorageUploadRequest{
		Data:        []byte("test data"),
		ContentType: "image/jpeg",
		Filename:    "test.jpg",
		BucketName:  testBucket,
		ObjectKey:   testKey,
	}

	assert.Equal(t, []byte("test data"), req.Data)
	assert.Equal(t, "image/jpeg", req.ContentType)
	assert.Equal(t, "test.jpg", req.Filename)
	assert.Equal(t, testBucket, req.BucketName)
	assert.Equal(t, testKey, req.ObjectKey)
}

func TestStorageUploadResponseValidation(t *testing.T) {
	resp := &entity.StorageUploadResponse{
		URL:       "https://example.com/test.jpg",
		ObjectKey: testKey,
	}

	assert.Equal(t, "https://example.com/test.jpg", resp.URL)
	assert.Equal(t, testKey, resp.ObjectKey)
}

// Test GCS repository creation with invalid credentials
func TestNewGCSStorageRepositoryInvalidCredentials(t *testing.T) {
	logger := &mockLogger{}
	
	// This should panic with invalid credentials
	assert.Panics(t, func() {
		repository.NewGCSStorageRepository(logger, "invalid-json")
	})
}

// Mock storage repository for testing interface compliance
type mockStorageRepository struct{}

func (m *mockStorageRepository) Upload(ctx context.Context, req *entity.StorageUploadRequest) (*entity.StorageUploadResponse, error) {
	return &entity.StorageUploadResponse{
		URL:       "https://example.com/" + req.ObjectKey,
		ObjectKey: req.ObjectKey,
	}, nil
}

func (m *mockStorageRepository) Download(ctx context.Context, bucketName, objectKey string) ([]byte, error) {
	return []byte("mock file data"), nil
}

func (m *mockStorageRepository) Delete(ctx context.Context, bucketName, objectKey string) error {
	return nil
}

func (m *mockStorageRepository) GetSignedURL(ctx context.Context, bucketName, objectKey string, expiration time.Duration) (string, error) {
	return "https://signed-url.example.com/" + objectKey, nil
}

func (m *mockStorageRepository) Close() error {
	return nil
}

func TestMockStorageRepositoryInterfaceCompliance(t *testing.T) {
	var repo repository.StorageRepository = &mockStorageRepository{}

	ctx := context.Background()

	// Test Upload
	req := &entity.StorageUploadRequest{
		Data:        []byte("test data"),
		ContentType: "image/jpeg",
		Filename:    "test.jpg",
		BucketName:  testBucket,
		ObjectKey:   testKey,
	}
	resp, err := repo.Upload(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, testKey, resp.ObjectKey)
	assert.Equal(t, "https://example.com/"+testKey, resp.URL)

	// Test Download
	data, err := repo.Download(ctx, testBucket, testKey)
	assert.NoError(t, err)
	assert.Equal(t, []byte("mock file data"), data)

	// Test Delete
	err = repo.Delete(ctx, testBucket, testKey)
	assert.NoError(t, err)

	// Test GetSignedURL
	url, err := repo.GetSignedURL(ctx, testBucket, testKey, time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, "https://signed-url.example.com/"+testKey, url)

	// Test Close
	err = repo.Close()
	assert.NoError(t, err)
}
