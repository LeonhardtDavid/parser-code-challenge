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
	"time"
)

type Crawler struct {
	maxParallelism uint
	scanner        *scanner.Scanner
	storage        storage.Storage

	timeout time.Duration
	visited map[string]bool
	mu      sync.Mutex
}

func (c *Crawler) markAsVisited(url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.visited[removeTrailingSlash(url)] = true
}

func (c *Crawler) wasVisited(url string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
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

func (c *Crawler) RecursiveScanAndSave(ctx context.Context, url *netUrl.URL) error {
	maxWorkers := c.maxParallelism
	var wg sync.WaitGroup
	urls := make(chan *netUrl.URL, maxWorkers*2)

	for i := uint(0); i < maxWorkers; i++ {
		wg.Add(1)
		go c.worker(ctx, &wg, urls, url.Host)
	}

	urls <- url

	wg.Wait()
	close(urls)

	return nil
}

func (c *Crawler) worker(ctx context.Context, wg *sync.WaitGroup, urls chan *netUrl.URL, host string) {
	//defer wg.Done()
	go c.checkDone(wg, urls) // TODO is there a better way to check if there is no more links? this is a clumsy solution
	for url := range urls {
		urlString := url.String()
		if !c.wasVisited(urlString) {
			c.markAsVisited(urlString)
			if result, err := c.ScanAndStore(ctx, url); err == nil {
				for _, link := range result.Links {
					if parsedLink, err := netUrl.ParseRequestURI(link); err == nil && parsedLink.Host == host {
						go func() { // Running in a goroutine to avoid blocks in case we max out the capacity of the channel (links get queued faster that processed)
							urls <- parsedLink
						}()
					}
				}
			}
		}
	}
}

func (c *Crawler) checkDone(wg *sync.WaitGroup, urls chan *netUrl.URL) {
	for {
		time.Sleep(c.timeout)
		if len(urls) == 0 {
			wg.Done()
			return
		}
	}
}

type Options = func(crawler *Crawler)

func WithMaxParallelism(parallelism uint) Options {
	return func(s *Crawler) {
		s.maxParallelism = parallelism
	}
}

func WithTimeout(timeout time.Duration) Options {
	return func(s *Crawler) {
		s.timeout = timeout
	}
}

func NewCrawler(scanner *scanner.Scanner, storage storage.Storage, options ...Options) *Crawler {
	c := &Crawler{
		maxParallelism: 1,
		scanner:        scanner,
		storage:        storage,
		timeout:        5 * time.Second,
		visited:        make(map[string]bool),
	}

	for _, o := range options {
		o(c)
	}

	return c
}
