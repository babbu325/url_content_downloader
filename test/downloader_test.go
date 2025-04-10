package test

import (
	model "jrog_assignment/internal/models"
	"jrog_assignment/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloaderMock(t *testing.T) {
	mockDownloader := new(mocks.Downloader)
	job := model.DownloadJob{URL: "http://example.com"}

	mockDownloader.On("Download", job).Return(model.DownloadResult{
		URL:     job.URL,
		Content: []byte("mock content"),
		Err:     nil,
	})

	result := mockDownloader.Download(job)
	assert.Nil(t, result.Err)
	assert.Equal(t, []byte("mock content"), result.Content)
	assert.Equal(t, "http://example.com", result.URL)
}
