package test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	model "jrog_assignment/internal/models"
	"jrog_assignment/mocks"
	"testing"
)

func TestReaderMock(t *testing.T) {
	mockReader := new(mocks.Reader)
	out := make(chan model.DownloadJob, 1)

	mockReader.On("Read", "dummy.csv", mock.AnythingOfType("chan<- model.DownloadJob")).Run(func(args mock.Arguments) {
		ch := args.Get(1).(chan<- model.DownloadJob)
		ch <- model.DownloadJob{URL: "http://example.com"}
		close(ch)
	}).Return(nil)

	err := mockReader.Read("dummy.csv", out)
	assert.Nil(t, err)

	job := <-out
	assert.Equal(t, "http://example.com", job.URL)
}
