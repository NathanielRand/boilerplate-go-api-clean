package handlers

import (
	"net/http"
	"sync"
)

// ResponseWriterWrapper is a wrapper around http.ResponseWriter that
// provides mutex locking around the WriteHeader and Write methods.
// This is necessary because the image conversion process is asynchronous
// and the response writer is shared between the request handler and the
// image conversion process.
type ResponseWriterWrapper struct {
	w  http.ResponseWriter
	mu sync.Mutex
}

// Write writes the data to the connection as part of an HTTP reply.
func (rw *ResponseWriterWrapper) Write(p []byte) (n int, err error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	return rw.w.Write(p)
}

// Header returns the header map that will be sent by WriteHeader.
func (rw *ResponseWriterWrapper) Header() http.Header {
	return rw.w.Header()
}

// WriteHeader sends an HTTP response header with status code.
func (rw *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	rw.w.WriteHeader(statusCode)
}
