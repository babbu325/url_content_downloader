package main

import (
	"flag"
	"fmt"
	"jrog_assignment/internal/app_pipeline"
)

func main() {
	filePath := flag.String("input", "urls.csv", "Path to input CSV file")
	flag.Parse()
	fmt.Println("file path: %v", *filePath)

	pipeline := app_pipeline.GetPipeline()
	pipeline.RunPipeline(*filePath)
}
