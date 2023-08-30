package app

import (
	"context"
	"fmt"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/scanner"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/storage"
	netUrl "net/url"
)

type Crawler struct {
	url     *netUrl.URL
	rawUrl  string
	scanner *scanner.Scanner
	storage storage.Storage
}

func (c *Crawler) ScanAndStore(ctx context.Context) error {
	visitedPage, err := c.scanner.LookupForLinks(ctx, c.url)
	if err != nil {
		return err
	}

	if err := c.storage.Save(ctx, visitedPage); err != nil {
		return fmt.Errorf("error storing visited pages: %w", err)
	}

	return nil
}

func (c *Crawler) RecursiveScanAndSave(ctx context.Context) error {
	visited := make(map[string]bool)
	return c.recursiveScanAndSaveWithVisitedCheck(ctx, visited)
}

func (c *Crawler) recursiveScanAndSaveWithVisitedCheck(ctx context.Context, visited map[string]bool) error {
	// TODO check for trailing slash
	visited[c.rawUrl] = true

	result, err := c.scanner.LookupForLinks(ctx, c.url)
	if err != nil {
		return err
	}

	if err := c.storage.Save(ctx, result); err != nil {
		return fmt.Errorf("error storing visited page: %w", err)
	}

	for _, link := range result.Links {
		_, exists := visited[link]
		if parsedLink, err := netUrl.ParseRequestURI(link); err == nil && !exists && parsedLink.Host == c.url.Host {
			visited[link] = true
			crawler := NewCrawler(parsedLink, c.scanner, c.storage)
			err := crawler.recursiveScanAndSaveWithVisitedCheck(ctx, visited)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewCrawler(url *netUrl.URL, scanner *scanner.Scanner, storage storage.Storage) Crawler {
	return Crawler{
		url:     url,
		rawUrl:  url.String(),
		scanner: scanner,
		storage: storage,
	}
}
