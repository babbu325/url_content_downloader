package test

import (
	"github.com/stretchr/testify/assert"
	model "jrog_assignment/internal/models"
	"jrog_assignment/mocks"
	"testing"
)

func TestWriterMock(t *testing.T) {
	mockWriter := new(mocks.Writer)
	result := model.DownloadResult{
		URL:     "http://example.com",
		Content: []byte("sample data"),
	}

	mockWriter.On("Write", result).Return(nil)

	err := mockWriter.Write(result)
	assert.Nil(t, err)
}
