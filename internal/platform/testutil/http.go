package testutil

import (
	"io"
	"net/http"
	"net/url"
)

type MockHTTPResponseWriter interface {
	WriteHeader(statusCode int)
	Write([]byte) (int, error)
	Header() http.Header
	Body() []byte
}

type mockHTTPResponseWriter struct {
	header     http.Header
	statusCode int
	body       []byte
}

func (m *mockHTTPResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

func (m mockHTTPResponseWriter) Write(data []byte) (int, error) {
	m.body = append(m.body, data...)
	return len(data), nil
}

func (m mockHTTPResponseWriter) Header() http.Header {
	return m.header
}

func (m mockHTTPResponseWriter) Body() []byte {
	return m.body
}

func (m mockHTTPResponseWriter) StatusCode() int {
	return m.statusCode
}

func NewMockHTTPResponseWriter() MockHTTPResponseWriter {
	return &mockHTTPResponseWriter{
		header: make(http.Header),
	}
}

func NewMockHTTPRequest(method, urlString string, body io.ReadCloser) *http.Request {
	return &http.Request{
		Method: method,
		URL:    &url.URL{Path: urlString},
		Body:   body,
	}
}
