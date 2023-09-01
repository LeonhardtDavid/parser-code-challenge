package app

import (
	"context"
	"errors"
	"github.com/LeonhardtDavid/parser-code-challenge/tests/helper"
	"github.com/stretchr/testify/assert"
	netUrl "net/url"
	"testing"
)

const (
	testUrl = "http://test.com/test"
)

func Test_ScanAndStore_ErrorOnScanner(t *testing.T) {
	scanner := &helper.MockFailScanner{
		Error: errors.New("some scanner failure"),
	}
	storage := &helper.MockStorage{}

	crawler := NewCrawler(scanner, storage, WithMaxParallelism(2))

	url, _ := netUrl.Parse(testUrl)
	result, err := crawler.ScanAndStore(context.Background(), url)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, scanner.Error, err)
}

func Test_ScanAndStore_ErrorOnStorage(t *testing.T) {
	scanner := &helper.MockScanner{
		Links: map[string][]string{
			testUrl: {},
		},
	}
	storage := &helper.MockFailStorage{
		Error: errors.New("some storage failure"),
	}

	crawler := NewCrawler(scanner, storage, WithMaxParallelism(2))

	url, _ := netUrl.Parse(testUrl)
	result, err := crawler.ScanAndStore(context.Background(), url)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, storage.Error, errors.Unwrap(err))
}

func Test_ScanAndStoreAllLinks_NoLinks(t *testing.T) {
	scanner := &helper.MockScanner{
		Links: map[string][]string{
			testUrl: {},
		},
	}
	storage := &helper.MockStorage{}

	crawler := NewCrawler(scanner, storage, WithMaxParallelism(2))

	url, _ := netUrl.Parse(testUrl)
	err := crawler.RecursiveScanAndStore(context.Background(), url)

	assert.Nil(t, err)
	assert.EqualValues(t, storage.Calls.Load(), 1)
	assert.EqualValues(t, len(storage.Links), 1)
	assert.True(t, helper.Contains(
		storage.Links,
		testUrl,
	))
}

func Test_ScanAndStoreAllLinks_MultipleLinks(t *testing.T) {
	scanner := &helper.MockScanner{
		Links: map[string][]string{
			testUrl: {
				"http://test.com",
				"http://test.com/test2",
				"http://test.com/test4",
				"http://othertest.com",
			},
			"http://test.com/test2": {
				"http://test.com",
				"http://test.com/test2",
				"http://test.com/test3",
				"http://othertest.com",
			},
			"http://test.com/test3": {
				"http://test.com/",
				"http://test.com/test2",
				"http://test.com/test3",
				"http://test.com/test4",
				"http://othertest.com",
			},
			"http://test.com/test4": {
				"http://test.com/test5",
			},
			"http://test.com/test5": {
				"http://test.com/test",
			},
		},
	}
	storage := &helper.MockStorage{}

	crawler := NewCrawler(scanner, storage, WithMaxParallelism(2))

	url, _ := netUrl.Parse(testUrl)
	err := crawler.RecursiveScanAndStore(context.Background(), url)

	assert.Nil(t, err)
	assert.EqualValues(t, storage.Calls.Load(), 6)
	assert.EqualValues(t, len(storage.Links), 6)
	assert.True(t, helper.Contains(
		storage.Links,
		testUrl,
		"http://test.com/test",
		"http://test.com/test2",
		"http://test.com/test3",
		"http://test.com/test4",
		"http://test.com/test5",
	))
}
