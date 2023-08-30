package app

import (
	"context"
	"fmt"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/model"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/storage"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	netUrl "net/url"
	"slices"
	"strings"
)

var (
	avoidSchemes = []string{"javascript"}
)

type Crawler struct {
	url     *netUrl.URL
	rawUrl  string
	baseUrl string
	storage storage.Storage
}

func (c *Crawler) Scan(ctx context.Context) (*model.VisitedPage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.rawUrl, nil)
	if err != nil {
		return nil, err // TODO better errors?
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err // TODO better errors?
	}
	defer res.Body.Close()

	result, err := c.getLinks(res.Body)
	if err != nil {
		return nil, err // TODO better errors?
	}

	return result, nil
}

func (c *Crawler) ScanAndStore(ctx context.Context) error {
	visitedPage, err := c.Scan(ctx)
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
	visited[c.rawUrl] = true
	result, err := c.Scan(ctx)
	if err != nil {
		return err
	}

	if err := c.storage.Save(ctx, result); err != nil {
		return fmt.Errorf("error storing visited pages: %w", err)
	}

	for _, link := range result.Links {
		_, exists := visited[link]
		if parsedLink, err := netUrl.ParseRequestURI(link); err == nil && !exists && parsedLink.Host == c.url.Host {
			visited[link] = true
			crawler := NewCrawler(parsedLink, c.storage)
			err := crawler.recursiveScanAndSaveWithVisitedCheck(ctx, visited)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Crawler) getLinks(reader io.Reader) (*model.VisitedPage, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	visitedPage := model.VisitedPage{
		Url: c.rawUrl,
	}

	doc.Find("a").Each(func(_ int, selection *goquery.Selection) {
		if ref, exists := selection.Attr("href"); exists {
			// TODO avoid listing duplicated links
			// avoids invalid links
			trimmedRef := strings.TrimSpace(ref)
			if parsedRef, err := netUrl.ParseRequestURI(trimmedRef); err == nil && !slices.Contains(avoidSchemes, parsedRef.Scheme) {
				var link string

				if strings.HasPrefix(trimmedRef, "//") {
					link = fmt.Sprintf("%s:%s", c.url.Scheme, trimmedRef)
				} else if strings.HasPrefix(trimmedRef, "/") {
					link = c.baseUrl + trimmedRef
				} else {
					link = trimmedRef
				}

				visitedPage.Links = append(visitedPage.Links, link)
			}
		}
	})

	return &visitedPage, nil
}

func NewCrawler(url *netUrl.URL, storage storage.Storage) Crawler {
	return Crawler{
		url:     url,
		rawUrl:  url.String(),
		baseUrl: fmt.Sprintf("%s://%s", url.Scheme, url.Host),
		storage: storage,
	}
}
