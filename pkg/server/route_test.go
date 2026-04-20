package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_registerFileServerMiddleware(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		path string
		code int
	}{
		{"Register file server middleware - Found page", "/", 200},
		{"Register file server middleware - Not found page", "/try_to_accces_non_existed_path", 404},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer()
			s.registerFileServerMiddleware()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.path, nil)
			s.ginEngine.ServeHTTP(w, req)

			if tt.code != w.Code {
				t.Fatalf("w.Code should return %d for path %s, got status code `%d`", tt.code, tt.path, w.Code)
			}
		})
	}
}
