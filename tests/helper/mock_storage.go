package helper

import (
	"context"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/model"
	"sync"
	"sync/atomic"
)

var mu sync.Mutex

type MockStorage struct {
	Calls atomic.Int64
	Links []string
}

func (ms *MockStorage) Save(_ context.Context, visitedPage *model.VisitedPage) error {
	ms.Calls.Add(1)
	mu.Lock()
	defer mu.Unlock()
	ms.Links = append(ms.Links, visitedPage.Url)

	return nil
}

type MockFailStorage struct {
	Error error
}

func (ms *MockFailStorage) Save(_ context.Context, _ *model.VisitedPage) error {
	return ms.Error
}
