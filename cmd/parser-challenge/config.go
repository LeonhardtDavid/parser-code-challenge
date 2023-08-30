package main

import (
	"errors"
	"flag"
	"fmt"
	netUrl "net/url"
	"strings"
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
	if !strings.HasPrefix(uri.Scheme, "http") {
		return nil, errors.New("url must start with http or https")
	}

	return &Config{Url: uri}, nil
}
