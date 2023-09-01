package scanner

import (
	"context"
	"errors"
	"github.com/LeonhardtDavid/parser-code-challenge/tests/helper"
	"github.com/stretchr/testify/assert"
	"net/http"
	netUrl "net/url"
	"testing"
)

const (
	testUrl = "http://test.com/test"
)

func Test_ErrorOnHttpResponse(t *testing.T) {
	client := &http.Client{
		Transport: helper.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("failure expected")
		}),
	}
	scanner := New(WithClient(client))

	url, _ := netUrl.Parse(testUrl)
	result, err := scanner.LookupForLinks(context.Background(), url)

	assert.Nil(t, result)
	assert.NotNil(t, err)
}

func Test_SuccessResponse(t *testing.T) {
	client := &http.Client{
		Transport: helper.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return helper.NewResponse(
				req,
				`<html>
					<head>
					</head>
					<body>
					<a href="http://test.com">link</a>
					<a href="http://test.com/test">link</a>
					<a href="http://test.com/test2">link</a>
					<a href="#">link</a>
					<a href="javascript:void(0)">JS link</a>
					<div>
						<div>
							<a href="http://test.com/test3">link</a>
							<span><a href="http://test.com/test4">link</a></span>
							<span><a href="//test.com/test5">link</a></span>
							<span><a href="/test6">link</a></span>
							<span><a href="http://othertest.com">link</a></span>
						</div>
					</div>
					</body>
				</html>`,
			), nil
		}),
	}
	scanner := New(WithClient(client))

	url, _ := netUrl.Parse(testUrl)
	result, err := scanner.LookupForLinks(context.Background(), url)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.True(t, helper.Contains(
		result.Links,
		"http://test.com",
		"http://test.com/test",
		"http://test.com/test2",
		"http://test.com/test3",
		"http://test.com/test4",
		"http://test.com/test5",
		"http://test.com/test6",
		"http://othertest.com",
	))
}
