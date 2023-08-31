package main

import (
	"context"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/app"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/scanner"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/storage"
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Println("Starting crawler...")

	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error on set up: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	crawler := app.NewCrawler(
		scanner.New(),
		storage.NewStdoutStorage(),
		app.WithMaxParallelism(config.Parallelism),
		app.WithTimeout(config.Timeout),
	)

	if err := crawler.RecursiveScanAndSave(ctx, config.Url); err != nil {
		log.Fatalf("Error crawling %q with error %v", config.Url, err)
	}
}
