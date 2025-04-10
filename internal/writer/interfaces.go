package writer

import model "jrog_assignment/internal/models"

type Writer interface {
	Write(result model.DownloadResult) error
}
