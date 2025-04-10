package writer

import (
	"fmt"
	"github.com/google/uuid"
	model "jrog_assignment/internal/models"
	"os"
)

type FileWriter struct{}

func (r *FileWriter) Write(in model.DownloadResult) error {

	name := "File_" + uuid.New().String() + ".txt"

	file, err := os.Create(fmt.Sprintf("urlContent/%s", name))
	if err != nil {
		return err
	}

	_, err = file.Write(in.Content)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}
