package main

import (
	"github.com/LeonhardtDavid/parser-code-challenge/internal/app"
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Println("Starting crawler...")

	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error on set up: %v", err)
	}

	crawler := app.NewCrawler(config.Url)

	result, err := crawler.Crawl()
	if err != nil {
		log.Fatalf("Error crawling %q with error %v", config.Url, err)
	}

	log.Println(result)
}
