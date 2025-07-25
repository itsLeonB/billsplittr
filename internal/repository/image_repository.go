package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"google.golang.org/api/option"
)

// gcsRepository is the concrete implementation using Google Cloud Storage.
type gcsRepository struct {
	client         *storage.Client
	bucketName     string
	serviceAccount string
	privateKey     []byte
}

func NewImageRepository(bucketName string, serviceAccountKey string) ImageRepository {
	// Create a GCS client using the key file
	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(serviceAccountKey)))
	if err != nil {
		log.Fatalf("failed to initialize GCS client: %v", err)
	}

	var parsed struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
	}
	if err = json.Unmarshal([]byte(serviceAccountKey), &parsed); err != nil {
		log.Fatalf("failed to unmarshal credentials JSON: %v", err)
	}

	return &gcsRepository{
		client:         client,
		bucketName:     bucketName,
		serviceAccount: parsed.ClientEmail,
		privateKey:     []byte(parsed.PrivateKey),
	}
}

// Upload uploads an image blob to GCS and returns the object name.
func (gr *gcsRepository) Upload(ctx context.Context, reader io.Reader, contentType string) (string, error) {
	ext := strings.Split(contentType, "/")[1]
	objectName := fmt.Sprintf("uploads/bill_%s.%s", uuid.New().String(), ext)

	wc := gr.client.Bucket(gr.bucketName).Object(objectName).NewWriter(ctx)
	wc.ContentType = contentType
	wc.CacheControl = "no-cache"

	if _, err := io.Copy(wc, reader); err != nil {
		return "", eris.Wrap(err, "failed to stream file to GCS")
	}

	if err := wc.Close(); err != nil {
		return "", eris.Wrap(err, "failed to close GCS writer")
	}

	return objectName, nil
}

// GenerateSignedURL returns a signed URL to access the uploaded object temporarily.
func (gr *gcsRepository) GenerateSignedURL(ctx context.Context, objectName string, duration time.Duration) (string, error) {
	url, err := storage.SignedURL(gr.bucketName, objectName, &storage.SignedURLOptions{
		Method:         "GET",
		Expires:        time.Now().Add(duration),
		GoogleAccessID: gr.serviceAccount,
		PrivateKey:     gr.privateKey,
		Scheme:         storage.SigningSchemeV4,
	})
	if err != nil {
		return "", eris.Wrap(err, "failed to generate signed URL")
	}

	return url, nil
}

func (gr *gcsRepository) Delete(ctx context.Context, objectName string) error {
	err := gr.client.Bucket(gr.bucketName).Object(objectName).Delete(ctx)
	if err != nil {
		return eris.Wrap(err, "failed to delete image from GCS")
	}
	return nil
}
