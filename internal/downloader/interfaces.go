package downloader

import (
	model "jrog_assignment/internal/models"
)

type Downloader interface {
	Download(in model.DownloadJob) model.DownloadResult
}
