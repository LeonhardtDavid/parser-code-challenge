package app

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
)

type Crawler struct {
	url string
}

func (c *Crawler) Crawl() ([]string, error) {
	res, err := http.Get(c.url)
	if err != nil {
		return nil, err // TODO better errors?
	}

	result, err := getLinks(res.Body)
	if err != nil {
		return nil, err // TODO better errors?
	}

	return result, nil
}

func getLinks(reader io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var links []string
	doc.Find("a").Each(func(_ int, selection *goquery.Selection) {
		if ref, exists := selection.Attr("href"); exists {
			links = append(links, ref)
		}
	})

	return links, nil
}

func NewCrawler(url string) Crawler {
	return Crawler{url: url}
}
