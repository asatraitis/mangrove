package utils

import (
	"io"
	"net/http"
)

type MockResponseWriter struct {
	Body io.ReadCloser
	Code int
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{}
}

func (w *MockResponseWriter) Header() http.Header {
	return http.Header{}
}
func (w *MockResponseWriter) Write(b []byte) (int, error) {
	return 0, nil
}
func (w *MockResponseWriter) WriteHeader(code int) {
	w.Code = code
}
