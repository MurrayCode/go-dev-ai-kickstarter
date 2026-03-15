package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "get returns hello world",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			wantBody:   "hello, world\n",
		},
		{
			name:       "post not allowed",
			method:     http.MethodPost,
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/hello", nil)
			rr := httptest.NewRecorder()

			helloHandler(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("helloHandler() status = %d, want %d", rr.Code, tc.wantStatus)
			}

			if rr.Body.String() != tc.wantBody {
				t.Fatalf("helloHandler() body = %q, want %q", rr.Body.String(), tc.wantBody)
			}
		})
	}
}

func TestNewMuxRoutesHelloEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rr := httptest.NewRecorder()

	NewMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("NewMux() /hello status = %d, want %d", rr.Code, http.StatusOK)
	}
}
