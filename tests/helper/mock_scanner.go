package helper

import (
	"context"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/model"
	netUrl "net/url"
)

type MockScanner struct {
	Links map[string][]string
}

func (ms *MockScanner) LookupForLinks(_ context.Context, url *netUrl.URL) (*model.VisitedPage, error) {
	if links, exists := ms.Links[url.String()]; exists {
		return &model.VisitedPage{
			Url:   url.String(),
			Links: links,
		}, nil
	}

	return &model.VisitedPage{
		Url:   url.String(),
		Links: []string{},
	}, nil
}

type MockFailScanner struct {
	Error error
}

func (ms *MockFailScanner) LookupForLinks(_ context.Context, url *netUrl.URL) (*model.VisitedPage, error) {
	return nil, ms.Error
}
