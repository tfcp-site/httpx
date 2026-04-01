package correlation_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tfcp-site/httpx/correlation"
)

func TestMiddleware_generatesIDWhenAbsent(t *testing.T) {
	var gotID string
	h := correlation.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID = correlation.FromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if gotID == "" {
		t.Fatal("expected a generated ID, got empty string")
	}
	if got := rr.Header().Get(correlation.Header); got != gotID {
		t.Fatalf("response header %q = %q, want %q", correlation.Header, got, gotID)
	}
}

func TestMiddleware_propagatesIncomingID(t *testing.T) {
	const wantID = "incoming-id-123"
	var gotID string
	h := correlation.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotID = correlation.FromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(correlation.Header, wantID)
	h.ServeHTTP(httptest.NewRecorder(), req)

	if gotID != wantID {
		t.Fatalf("FromContext = %q, want %q", gotID, wantID)
	}
}

func TestTransport_injectsIDFromContext(t *testing.T) {
	const wantID = "outgoing-id-456"
	var gotHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get(correlation.Header)
	}))
	defer server.Close()

	client := &http.Client{
		Transport: &correlation.Transport{Base: http.DefaultTransport},
	}

	ctx := correlation.WithContext(context.Background(), wantID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	_, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if gotHeader != wantID {
		t.Fatalf("server got %q = %q, want %q", correlation.Header, gotHeader, wantID)
	}
}

func TestTransport_skipsWhenContextEmpty(t *testing.T) {
	var gotHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get(correlation.Header)
	}))
	defer server.Close()

	client := &http.Client{
		Transport: &correlation.Transport{Base: http.DefaultTransport},
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	_, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if gotHeader != "" {
		t.Fatalf("expected no header, got %q", gotHeader)
	}
}
