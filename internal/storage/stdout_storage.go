package storage

import (
	"context"
	"encoding/json"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/model"
	"log"
)

type stdoutStorage struct{}

func (s *stdoutStorage) SaveAll(_ context.Context, visitedPages []model.VisitedPage) error {
	js, err := json.MarshalIndent(visitedPages, "", "  ")
	if err == nil {
		log.Println(string(js))
	}

	return err
}

func NewStdoutStorage() Storage {
	return &stdoutStorage{}
}
