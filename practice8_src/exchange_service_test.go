package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestServer(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}))
}

func TestGetRate_Success(t *testing.T) {
	srv := newTestServer(http.StatusOK, `{"base":"USD","target":"EUR","rate":0.92}`)
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	rate, err := svc.GetRate("USD", "EUR")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if rate != 0.92 {
		t.Errorf("expected rate 0.92, got %v", rate)
	}
}

func TestGetRate_APIBusinessError_404(t *testing.T) {
	srv := newTestServer(http.StatusNotFound, `{"error":"invalid currency pair"}`)
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "XYZ")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "api error: invalid currency pair"
	if err.Error() != want {
		t.Errorf("expected %q, got %q", want, err.Error())
	}
}

func TestGetRate_APIBusinessError_400(t *testing.T) {
	srv := newTestServer(http.StatusBadRequest, `{"error":"invalid currency pair"}`)
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("", "EUR")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	want := "api error: invalid currency pair"
	if err.Error() != want {
		t.Errorf("expected %q, got %q", want, err.Error())
	}
}

func TestGetRate_MalformedJSON(t *testing.T) {
	srv := newTestServer(http.StatusOK, `Internal Server Error`)
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected decode error, got nil")
	}
	if len(err.Error()) == 0 {
		t.Error("expected non-empty error message")
	}
}

func TestGetRate_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// sleep longer than client timeout
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	svc.Client = &http.Client{Timeout: 50 * time.Millisecond}

	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestGetRate_ServerPanic500(t *testing.T) {
	srv := newTestServer(http.StatusInternalServerError, `{"error":"internal server error"}`)
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected error for 500, got nil")
	}
}

func TestGetRate_EmptyBody(t *testing.T) {
	srv := newTestServer(http.StatusOK, ``)
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "EUR")
	if err == nil {
		t.Fatal("expected error for empty body, got nil")
	}
}
