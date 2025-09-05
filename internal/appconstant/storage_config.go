package appconstant

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB

	StatusPending   = "pending"
	StatusUploaded  = "uploaded"
	StatusProcessed = "processed"
	StatusFailed    = "failed"
)

var AllowedContentTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/jpg":  {},
	"image/png":  {},
	"image/webp": {},
}
