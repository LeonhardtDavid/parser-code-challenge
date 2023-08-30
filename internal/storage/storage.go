package storage

import (
	"context"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/model"
)

type Storage interface {
	Save(ctx context.Context, visitedPage *model.VisitedPage) error
}
