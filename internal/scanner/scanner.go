package scanner

import (
	"context"
	"fmt"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/model"
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

type Scanner interface {
	LookupForLinks(ctx context.Context, url *netUrl.URL) (*model.VisitedPage, error)
}

type scannerImpl struct {
	httpClient *http.Client
}

func (s scannerImpl) LookupForLinks(ctx context.Context, url *netUrl.URL) (*model.VisitedPage, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err // TODO better errors?
	}
	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err // TODO better errors?
	}
	defer res.Body.Close()

	result, err := getLinks(url, res.Body)
	if err != nil {
		return nil, err // TODO better errors?
	}

	return result, nil
}

func getLinks(url *netUrl.URL, reader io.Reader) (*model.VisitedPage, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	visitedPage := model.VisitedPage{
		Url: url.String(),
	}

	baseUrl := fmt.Sprintf("%s://%s", url.Scheme, url.Host)

	doc.Find("a").Each(func(_ int, selection *goquery.Selection) {
		if ref, exists := selection.Attr("href"); exists {
			// TODO avoid listing duplicated links?
			// avoids invalid links
			trimmedRef := strings.TrimSpace(ref)
			if parsedRef, err := netUrl.ParseRequestURI(trimmedRef); err == nil && !slices.Contains(avoidSchemes, parsedRef.Scheme) {
				var link string

				if strings.HasPrefix(trimmedRef, "//") {
					link = fmt.Sprintf("%s:%s", url.Scheme, trimmedRef)
				} else if strings.HasPrefix(trimmedRef, "/") {
					link = baseUrl + trimmedRef
				} else {
					link = trimmedRef
				}

				visitedPage.Links = append(visitedPage.Links, link)
			}
		}
	})

	return &visitedPage, nil
}

type Options func(scanner Scanner)

func WithClient(client *http.Client) Options {
	return func(s Scanner) {
		switch impl := s.(type) {
		case *scannerImpl:
			impl.httpClient = client
		}
	}
}

func New(options ...Options) Scanner {
	s := &scannerImpl{
		httpClient: http.DefaultClient,
	}

	for _, o := range options {
		o(s)
	}

	return s
}
