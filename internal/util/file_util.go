package util

import (
	"mime/multipart"

	"github.com/itsLeonB/billsplittr/internal/appconstant"
)

func IsImageType(fileHeader *multipart.FileHeader) (string, bool) {
	contentType := fileHeader.Header.Get("Content-Type")
	_, ok := appconstant.ImageTypes[contentType]
	return contentType, ok
}
