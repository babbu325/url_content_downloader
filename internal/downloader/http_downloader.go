package downloader

import (
	"io"
	model "jrog_assignment/internal/models"
	"net/http"
	"time"
)

type HttpDownloader struct{}

func (d *HttpDownloader) Download(in model.DownloadJob) model.DownloadResult {
	start := time.Now()
	resp, err := http.Get(in.URL)
	if err != nil {
		return model.DownloadResult{URL: in.URL, Err: err}
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	_ = time.Since(start)

	return model.DownloadResult{
		URL:     in.URL,
		Content: body,
		Err:     err,
	}
}
