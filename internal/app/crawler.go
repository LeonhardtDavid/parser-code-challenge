package app

import (
	"context"
	"fmt"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/model"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/scanner"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/storage"
	netUrl "net/url"
	"sync"
)

type Crawler struct {
	maxParallelism uint
	scanner        *scanner.Scanner
	storage        storage.Storage

	visited map[string]bool
	mu      sync.Mutex
}

func (c *Crawler) markAsVisited(url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.visited[url] = true
}

func (c *Crawler) wasVisited(url string) bool {
	// TODO trailing slashes could make a same page get visited twice
	_, exists := c.visited[url]
	return exists
}

func (c *Crawler) ScanAndStore(ctx context.Context, url *netUrl.URL) (*model.VisitedPage, error) {
	visitedPage, err := c.scanner.LookupForLinks(ctx, url)
	if err != nil {
		return nil, err
	}

	if err := c.storage.Save(ctx, visitedPage); err != nil {
		return nil, fmt.Errorf("error storing visited pages: %w", err)
	}

	return visitedPage, nil
}

func (c *Crawler) RecursiveScanAndSave(ctx context.Context, url *netUrl.URL) error {
	maxWorkers := c.maxParallelism
	var wg sync.WaitGroup
	urls := make(chan *netUrl.URL, maxWorkers*2)

	for i := 0; i < int(maxWorkers); i++ {
		wg.Add(1)
		go c.worker(ctx, &wg, urls, url.Host)
	}

	urls <- url

	wg.Wait()
	close(urls)

	return nil
}

func (c *Crawler) worker(ctx context.Context, wg *sync.WaitGroup, urls chan *netUrl.URL, host string) {
	defer wg.Done()
	for url := range urls {
		c.markAsVisited(url.String())
		if result, err := c.ScanAndStore(ctx, url); err == nil {
			for _, link := range result.Links {
				if parsedLink, err := netUrl.ParseRequestURI(link); err == nil && parsedLink.Host == host && !c.wasVisited(link) {
					urls <- parsedLink
				}
			}
		}
	}
}

type Options = func(crawler *Crawler)

func WithMaxParallelism(parallelism uint) Options {
	return func(s *Crawler) {
		s.maxParallelism = parallelism
	}
}

func NewCrawler(scanner *scanner.Scanner, storage storage.Storage, options ...Options) *Crawler {
	c := &Crawler{
		maxParallelism: 1,
		scanner:        scanner,
		storage:        storage,
		visited:        make(map[string]bool),
	}

	for _, o := range options {
		o(c)
	}

	return c
}
