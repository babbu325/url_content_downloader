package main

import (
	"jrog_assignment/internal/app_pipeline"
)

func main() {
	pipeline := app_pipeline.GetPipeline()
	pipeline.RunPipeline()
}
