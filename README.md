# Concurrent CSV URL Downloader (Go)

##  Overview
This is a Go command-line application that reads a CSV file containing a list of URLs,
downloads each URL’s content using up to 50 goroutines, and saves the content to `.txt` files.

The system is designed as a **3-stage concurrent pipeline**:
1. **Reader**: Streams URLs from the CSV file.
2. **Downloader**: Downloads content in parallel using up to 50 goroutines.
3. **Writer**: Saves each downloaded result to a file.

---

## Installation & Running

###  Prerequisites
- Go 1.23 or later installed
- A `csv` file with one URL per line (CSV format)

###  Setup
- git clone https://github.com/your/repo-name.git
- cd repo-name
- go mod tidy


###  Run the Application
- go run cmd/main.go -input "csx file path"

Output files will be saved to the project `urlContent/` directory as `.txt` file.
Make sure `urlContent` dir exists. 

---

## Running Unit Tests

### Install Mockery (for mocks)
```bash
go install github.com/vektra/mockery/v2@latest
```

Generate mocks:
```bash
mockery --all --recursive --output=mocks --case=camel
```

### Run All Tests
```bash
go test ./test/...
```
---

## Design Decisions

### Interfaces
- Defined `Reader`, `Downloader`, and `Writer` interfaces for better testability and separation of concerns.
- Used [`mockery`](https://github.com/vektra/mockery) to generate mocks for each component.

### Concurrency
- Used channels to communicate between stages.
- Capped downloader goroutines at 50 to control concurrency and resource usage.

### Graceful Shutdown
- On `SIGINT` or `SIGTERM`, the app allows a 5-second grace period to finish in-flight tasks before shutting down.

---

## Dependencies
 
mockery:       Mock generation (testing) 
stretchr/testify:  Assertions in tests 
