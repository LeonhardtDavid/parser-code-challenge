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
	if uri, err := netUrl.ParseRequestURI(url); err != nil {
		return nil, fmt.Errorf("the provided url is invalid: %w", err)
	} else if uri.Path != "" && uri.Path != "/" {
		return nil, errors.New("the url should not contain a path")
	} else if uri.RawQuery != "" {
		return nil, errors.New("the url should not contain any query parameter")
	}

	return &Config{Url: url}, nil
}
