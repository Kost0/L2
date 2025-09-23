package main

import (
	"L2/16/downloader"
	"flag"
	"log"
	"os"
	"time"
)

func main() {
	url := flag.String("url", "", "URL to mirror")
	depth := flag.Int("depth", 3, "Recursion depth")
	outputDir := flag.String("output", "mirrored_site", "Output directory")
	concurrency := flag.Int("concurrency", 5, "Number of concurrent mirrors")
	timeout := flag.Int("timeout", 30, "Request timeout")

	flag.Parse()

	if *url == "" {
		log.Fatal("Missing URL")
	}

	err := os.MkdirAll(*outputDir, 0755)
	if err != nil {
		log.Fatal("Failed to create output directory")
	}

	dl := downloader.NewDownloader(&downloader.Config{
		BaseURL:       *url,
		OutputDir:     *outputDir,
		MaxDepth:      *depth,
		Concurrency:   *concurrency,
		Timeout:       time.Duration(*timeout) * time.Second,
		UserAgent:     "WebMirror/1.0",
		RespectRobots: true,
	})

	err = dl.Start()
	if err != nil {
		log.Fatal("Mirroring failed: %v", err)
	}
}
