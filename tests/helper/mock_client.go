package helper

import (
	"bytes"
	"io"
	"net/http"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func NewResponse(req *http.Request, bodyString string) *http.Response {
	var body io.ReadCloser
	if bodyString == "" {
		body = io.NopCloser(bytes.NewReader(nil))
	} else {
		body = io.NopCloser(bytes.NewBufferString(bodyString))
	}

	return &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		Body:          body,
		ContentLength: int64(len(bodyString)),
		Request:       req,
	}
}
