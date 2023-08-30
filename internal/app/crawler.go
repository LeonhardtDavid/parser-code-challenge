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
	"strings"
)

type Crawler struct {
	url     string
	storage storage.Storage
}

func (c *Crawler) Scan(_ context.Context) ([]model.VisitedPage, error) {
	res, err := http.Get(c.url) // TODO use context
	if err != nil {
		return nil, err // TODO better errors?
	}
	defer res.Body.Close()

	result, err := c.getLinks(res.Body)
	if err != nil {
		return nil, err // TODO better errors?
	}

	return []model.VisitedPage{*result}, nil
}

func (c *Crawler) ScanAndStore(ctx context.Context) error {
	visitedPages, err := c.Scan(ctx)
	if err != nil {
		return err
	}

	if err := c.storage.SaveAll(ctx, visitedPages); err != nil {
		return fmt.Errorf("error storing visited pages: %w", err)
	}

	return nil
}

func (c *Crawler) getLinks(reader io.Reader) (*model.VisitedPage, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	visitedPage := model.VisitedPage{
		Url: c.url,
	}

	doc.Find("a").Each(func(_ int, selection *goquery.Selection) {
		if ref, exists := selection.Attr("href"); exists {
			// avoids invalid links
			trimmedRef := strings.TrimSpace(ref)
			if _, err := netUrl.ParseRequestURI(trimmedRef); err == nil {
				var link string

				if strings.HasPrefix(trimmedRef, "/") {
					link = c.url + trimmedRef
				} else {
					link = trimmedRef
				}

				visitedPage.Links = append(visitedPage.Links, link)
			}
		}
	})

	return &visitedPage, nil
}

func NewCrawler(url string, storage storage.Storage) Crawler {
	return Crawler{
		url:     url,
		storage: storage,
	}
}
