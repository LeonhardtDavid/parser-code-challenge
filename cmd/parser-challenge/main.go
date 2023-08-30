package main

import (
	"context"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/app"
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

	crawler := app.NewCrawler(config.Url, storage.NewStdoutStorage())

	if err := crawler.ScanAndStore(ctx); err != nil {
		log.Fatalf("Error crawling %q with error %v", config.Url, err)
	}
}
