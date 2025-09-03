package util_test

import (
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/itsLeonB/billsplittr/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestFileUtil_IsImageType_ValidImage(t *testing.T) {
	header := &multipart.FileHeader{
		Header: textproto.MIMEHeader{
			"Content-Type": []string{"image/jpeg"},
		},
	}
	
	contentType, isImage := util.IsImageType(header)
	
	assert.Equal(t, "image/jpeg", contentType)
	assert.True(t, isImage)
}

func TestFileUtil_IsImageType_InvalidType(t *testing.T) {
	header := &multipart.FileHeader{
		Header: textproto.MIMEHeader{
			"Content-Type": []string{"text/plain"},
		},
	}
	
	contentType, isImage := util.IsImageType(header)
	
	assert.Equal(t, "text/plain", contentType)
	assert.False(t, isImage)
}

func TestFileUtil_IsImageType_NoContentType(t *testing.T) {
	header := &multipart.FileHeader{
		Header: textproto.MIMEHeader{},
	}
	
	contentType, isImage := util.IsImageType(header)
	
	assert.Equal(t, "", contentType)
	assert.False(t, isImage)
}
