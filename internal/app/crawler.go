package app

import (
	"context"
	"fmt"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/model"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/scanner"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/storage"
	netUrl "net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Crawler struct {
	maxParallelism uint
	scanner        scanner.Scanner
	storage        storage.Storage

	visited      map[string]bool
	visitCounter atomic.Int64
	mu           sync.RWMutex
}

func (c *Crawler) markAsVisited(url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.visited[removeTrailingSlash(url)] = true
}

func (c *Crawler) wasVisited(url string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.visited[removeTrailingSlash(url)]
	return exists
}

func removeTrailingSlash(url string) string {
	return strings.TrimSuffix(url, "/")
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

func (c *Crawler) RecursiveScanAndStore(ctx context.Context, url *netUrl.URL) error {
	maxWorkers := c.maxParallelism
	urls := make(chan *netUrl.URL, maxWorkers*2)

	for i := uint(0); i < maxWorkers; i++ {
		go c.worker(ctx, url.Host, urls)
	}

	c.visitCounter.Add(1)
	urls <- url

	for { // TODO is there a better way to check if there is no more links? this seems a clumsy solution
		time.Sleep(500 * time.Millisecond)
		if c.visitCounter.Load() == 0 && len(urls) == 0 {
			close(urls)
			break
		}
	}

	return nil
}

func (c *Crawler) worker(ctx context.Context, host string, urls chan *netUrl.URL) {
	for url := range urls {
		urlString := url.String()
		if !c.wasVisited(urlString) {
			c.markAsVisited(urlString)
			if result, err := c.ScanAndStore(ctx, url); err == nil {
				for _, link := range result.Links {
					if parsedLink, err := netUrl.ParseRequestURI(link); err == nil && parsedLink.Host == host {
						go func() { // Running in a goroutine to avoid blocks in case we max out the capacity of the channel (links get queued faster that processed)
							c.visitCounter.Add(1)
							urls <- parsedLink
						}()
					}
				}
			}
		}
		c.visitCounter.Add(-1)
	}
}

type Options = func(crawler *Crawler)

func WithMaxParallelism(parallelism uint) Options {
	return func(s *Crawler) {
		s.maxParallelism = parallelism
	}
}

func NewCrawler(scanner scanner.Scanner, storage storage.Storage, options ...Options) *Crawler {
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
