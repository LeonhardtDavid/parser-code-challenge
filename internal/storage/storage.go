package storage

import (
	"context"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/model"
)

type Storage interface {
	SaveAll(ctx context.Context, visitedPages []model.VisitedPage) error
}
