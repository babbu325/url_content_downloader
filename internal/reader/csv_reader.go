package reader

import (
	"bufio"
	"encoding/csv"
	"io"
	model "jrog_assignment/internal/models"

	"os"
)

type CsvReader struct{}

func (r *CsvReader) Read(filePath string, out chan<- model.DownloadJob) error {
	file, err := os.Open(filePath)
	if err != nil {

		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	//Read the first line of csv file and skip it as it's a header.
	if _, err := reader.Read(); err != nil {

	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) == 0 {
			continue
		}
		out <- model.DownloadJob{URL: record[0]}
	}
	
	return nil
}
