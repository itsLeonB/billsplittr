package entity

type StorageUploadRequest struct {
	Data        []byte
	ContentType string
	Filename    string
	BucketName  string
	ObjectKey   string
}

type StorageUploadResponse struct {
	URL       string
	ObjectKey string
}
