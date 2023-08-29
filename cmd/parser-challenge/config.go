package main

import (
	"errors"
	"flag"
	"fmt"
	netUrl "net/url"
)

type Config struct {
	Url string
}

func LoadConfig() (*Config, error) {
	var url string
	flag.StringVar(&url, "url", "", "url on which the crawler will run")

	flag.Parse()

	if url == "" {
		return nil, errors.New("url flag is required")
	}
	if _, err := netUrl.ParseRequestURI(url); err != nil {
		return nil, fmt.Errorf("the provided url is invalid: %w", err)
	}

	return &Config{Url: url}, nil
}
