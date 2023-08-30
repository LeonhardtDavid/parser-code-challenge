package main

import (
	"errors"
	"flag"
	"fmt"
	netUrl "net/url"
)

type Config struct {
	Url *netUrl.URL
}

func LoadConfig() (*Config, error) {
	var url string
	flag.StringVar(&url, "url", "", "url on which the crawler will run")

	flag.Parse()

	if url == "" {
		return nil, errors.New("url flag is required")
	}
	uri, err := netUrl.ParseRequestURI(url)
	if err != nil {
		return nil, fmt.Errorf("the provided url is invalid: %w", err)
	}

	return &Config{Url: uri}, nil
}
