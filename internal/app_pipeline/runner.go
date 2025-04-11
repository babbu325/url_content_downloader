package app_pipeline

import (
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

func (s *Stages) RunPipeline(filePath string) {

	urlChan := make(chan model.DownloadJob, 1000)
	resultChan := make(chan model.DownloadResult, 1000)
	done := make(chan bool)
	var wg sync.WaitGroup

	s.startReader(filePath, urlChan)
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
	const maxGoroutines = 50
	sem := make(chan struct{}, maxGoroutines)

	for job := range in {
		sem <- struct{}{}
		wg.Add(1)

		go func(j model.DownloadJob) {
			defer func() {
				<-sem
				wg.Done()
			}()
			result := s.Downloader.Download(j)
			out <- result
		}(job)
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
