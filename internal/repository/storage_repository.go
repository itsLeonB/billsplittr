package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/itsLeonB/billsplittr/internal/entity"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
	"google.golang.org/api/option"
)

type gcsStorageRepository struct {
	logger ezutil.Logger
	client *storage.Client
}

func NewGCSStorageRepository(logger ezutil.Logger, credentialsFile string) StorageRepository {

	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(credentialsFile)))
	if err != nil {
		panic(fmt.Sprintf("failed to create GCS client: %v", err))
	}

	return &gcsStorageRepository{logger, client}
}

func (r *gcsStorageRepository) Upload(ctx context.Context, req *entity.StorageUploadRequest) (*entity.StorageUploadResponse, error) {
	bucket := r.client.Bucket(req.BucketName)
	obj := bucket.Object(req.ObjectKey)

	// Create a writer to upload the file
	writer := obj.NewWriter(ctx)
	writer.ContentType = req.ContentType
	writer.Metadata = map[string]string{
		"original_filename": req.Filename,
		"uploaded_at":       time.Now().Format(time.RFC3339),
	}

	// Set cache control for images
	writer.CacheControl = "public, max-age=3600" // 1 hour cache

	// Write the file data
	if _, err := io.Copy(writer, bytes.NewReader(req.Data)); err != nil {
		_ = writer.Close() // best-effort close on copy failure
		return nil, eris.Wrap(err, "failed to upload file to GCS")
	}
	if err := writer.Close(); err != nil {
		return nil, eris.Wrap(err, "failed to finalize upload to GCS")
	}

	// // Make the object publicly readable (optional)
	// if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
	// 	// Log warning but don't fail the operation
	// 	// In production, you might want to handle this differently based on your security requirements
	// 	fmt.Printf("Warning: failed to make object public: %v\n", err)
	// }

	// // Generate public URL
	// publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", req.BucketName, req.ObjectKey)

	return &entity.StorageUploadResponse{
		// URL:       publicURL,
		ObjectKey: req.ObjectKey,
	}, nil
}

func (r *gcsStorageRepository) Download(ctx context.Context, bucketName, objectKey string) ([]byte, error) {
	bucket := r.client.Bucket(bucketName)
	obj := bucket.Object(objectKey)

	reader, err := obj.NewReader(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return nil, eris.New("file not found")
		}
		return nil, eris.Wrap(err, "failed to create GCS reader")
	}
	defer r.close(reader, "obj reader")

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, eris.Wrap(err, "failed to read file from GCS")
	}

	return data, nil
}

func (r *gcsStorageRepository) Delete(ctx context.Context, bucketName, objectKey string) error {
	bucket := r.client.Bucket(bucketName)
	obj := bucket.Object(objectKey)

	if err := obj.Delete(ctx); err != nil {
		if err == storage.ErrObjectNotExist {
			// Object doesn't exist, consider it already deleted
			return nil
		}
		return eris.Wrap(err, "failed to delete file from GCS")
	}

	return nil
}

func (r *gcsStorageRepository) GetSignedURL(ctx context.Context, bucketName, objectKey string, expiration time.Duration) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expiration),
	}

	url, err := storage.SignedURL(bucketName, objectKey, opts)
	if err != nil {
		return "", eris.Wrap(err, "failed to generate signed URL")
	}

	return url, nil
}

func (r *gcsStorageRepository) Close() error {
	return r.client.Close()
}

func (r *gcsStorageRepository) close(closeable io.Closer, objName string) {
	if err := closeable.Close(); err != nil {
		stkErr := eris.Wrapf(err, "failed to close %s", objName)
		r.logger.Errorf(eris.ToString(stkErr, true))
	}
}
