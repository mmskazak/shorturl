package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock handler to be wrapped by middleware.
func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message":"hello, world"}`))
}

func TestGzipMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		acceptEncoding string
		contentType    string
		expectGzip     bool
	}{
		{
			name:           "without gzip",
			acceptEncoding: "",
			contentType:    "application/json",
			expectGzip:     false,
		},
		{
			name:           "with non-gzip encoding",
			acceptEncoding: "deflate",
			contentType:    "application/json",
			expectGzip:     false,
		},
		{
			name:           "unsupported content type",
			acceptEncoding: "gzip",
			contentType:    "text/plain",
			expectGzip:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := bytes.NewBufferString(`{"data":"example"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", reqBody)
			req.Header.Add("Content-Type", tt.contentType)

			if tt.acceptEncoding != "" {
				req.Header.Add("Accept-Encoding", tt.acceptEncoding)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Wrap the handler with the middleware
			handler := GzipMiddleware(http.HandlerFunc(mockHandler))

			// Serve the HTTP request
			handler.ServeHTTP(rr, req)

			// Check the response
			result := rr.Result()
			defer func() {
				if err := result.Body.Close(); err != nil {
					log.Fatalf("failed to close response body: %v", err)
				}
			}()

			if tt.expectGzip {
				assert.Equal(t, "gzip", result.Header.Get("Content-Encoding"))

				gzipReader, err := gzip.NewReader(result.Body)
				require.NoError(t, err)
				defer func(gzipReader *gzip.Reader) {
					err := gzipReader.Close()
					if err != nil {
						log.Fatalf("error closing gzipReader: %v", err)
					}
				}(gzipReader)

				body, err := io.ReadAll(gzipReader)
				require.NoError(t, err)
				assert.JSONEq(t, `{"message":"hello, world"}`, string(body))
			} else {
				assert.NotEqual(t, "gzip", result.Header.Get("Content-Encoding"))

				body, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				assert.JSONEq(t, `{"message":"hello, world"}`, string(body))
			}
		})
	}
}
