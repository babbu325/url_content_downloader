package model

type DownloadJob struct {
	URL string
}

type DownloadResult struct {
	URL     string
	Content []byte
	Err     error
}
