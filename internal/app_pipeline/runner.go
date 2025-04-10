package app_pipeline

import (
	"flag"
	"fmt"
	"jrog_assignment/internal/downloader"
	"jrog_assignment/internal/logger"
	model "jrog_assignment/internal/models"
	"jrog_assignment/internal/reader"
	"jrog_assignment/internal/writer"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func init() {
	logger.InitLogger()
}

type Stages struct {
	Reader     reader.Reader
	Downloader downloader.Downloader
	Writer     writer.Writer
}

func GetPipeline() *Stages {
	return &Stages{
		Reader:     &reader.CsvReader{},
		Downloader: &downloader.HttpDownloader{},
		Writer:     &writer.FileWriter{},
	}

}

func (s *Stages) RunPipeline() {

	filePath := flag.String("input", "urls.csv", "Path to input CSV file")
	flag.Parse()
	fmt.Println("file path: %v", *filePath)

	urlChan := make(chan model.DownloadJob, 100)
	resultChan := make(chan model.DownloadResult, 100)
	done := make(chan bool)
	var wg sync.WaitGroup

	s.startReader(*filePath, urlChan)
	s.startDownloaders(urlChan, resultChan, &wg)
	s.startResultCloser(&wg, resultChan)
	s.startWriter(resultChan, done)
	s.handleShutdown(done)
}

func (s *Stages) startReader(filePath string, out chan<- model.DownloadJob) {
	go func() {

		err := s.Reader.Read(filePath, out)
		if err != nil {
			fmt.Println("Failed to read CSV:", err)
			close(out)
			return
		}
		close(out)
	}()
}

func (s *Stages) startDownloaders(in <-chan model.DownloadJob, out chan<- model.DownloadResult, wg *sync.WaitGroup) {
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range in {
				out <- s.Downloader.Download(job) // optional: collect for stats
			}
		}()
	}
}

func (s *Stages) startResultCloser(wg *sync.WaitGroup, out chan<- model.DownloadResult) {
	go func() {
		wg.Wait()
		close(out)
	}()

}

func (s *Stages) startWriter(in <-chan model.DownloadResult, done chan<- bool) {
	go func() {
		for result := range in {
			if result.Err != nil {
				fmt.Printf("Failed to download %s: %v", result.URL, result.Err)
				continue
			}
			if err := s.Writer.Write(result); err != nil {
				fmt.Println("Save failed:", err)
			}
		}
		done <- true
	}()
}

func (s *Stages) handleShutdown(done <-chan bool) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-shutdown:
		fmt.Println("Shutting down, 5s grace period...")
		time.Sleep(5 * time.Second)
	case <-done:
		fmt.Println("Completed all downloads.")
	}
}
