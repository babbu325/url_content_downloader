package reader

import model "jrog_assignment/internal/models"

type Reader interface {
	Read(filePath string, out chan<- model.DownloadJob) error
}
