package storage

import (
	"context"
	"encoding/json"
	"github.com/LeonhardtDavid/parser-code-challenge/internal/model"
	"log"
)

type stdoutStorage struct{}

func (s *stdoutStorage) Save(_ context.Context, visitedPage *model.VisitedPage) error {
	js, err := json.MarshalIndent(visitedPage, "", "  ")
	if err == nil {
		log.Println(string(js))
	}

	return err
}

func NewStdoutStorage() Storage {
	return &stdoutStorage{}
}
