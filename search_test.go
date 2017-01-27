package search

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		s, err := New()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := s.baseURL.String(), "https://search.b17g.services/"; got != want {
			t.Errorf("s.baseURL.String() = %q, want %q", got, want)
		}
	})

	t.Run("OptionReturningError", func(t *testing.T) {
		optionError := errors.New("option error")
		option := func(*Search) error {
			return optionError
		}

		_, err := New(SetBaseURL("/"), option)

		if got, want := err, optionError; got != want {
			t.Errorf("got err = %v, %v", got, want)
		}
	})

	t.Run("SetBaseURL", func(t *testing.T) {
		s, err := New(SetBaseURL("http://example.com/"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := s.baseURL.String(), "http://example.com/"; got != want {
			t.Errorf("s.baseURL.String() = %q, want %q", got, want)
		}
	})

	t.Run("SetLogf", func(t *testing.T) {
		var buf bytes.Buffer
		logf := func(format string, v ...interface{}) {
			fmt.Fprintf(&buf, format, v...)
		}
		s, err := New(SetBaseURL("/"), SetLogf(logf))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		s.logf("foo %s", "bar")

		if got, want := buf.String(), "foo bar"; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("SetHTTPClient", func(t *testing.T) {
		hc := &http.Client{}

		s, err := New(SetBaseURL("/"), SetHTTPClient(hc))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := s.httpClient, hc; got != want {
			t.Errorf("s.httpClient = %p, want %p", got, want)
		}
	})
}

func TestIsJSONResponse(t *testing.T) {
	for n, tt := range []struct {
		contentType string
		isJSON      bool
	}{
		{"application/json; charset=utf-8", true},
		{"application/json; charset=iso-8859-1", true},
		{"application/json", true},
		{"text/plain", false},
		{"randomnoiseapplication/jsonrandomnoise", false},
	} {
		resp := &http.Response{Header: make(http.Header)}
		resp.Header.Add("Content-Type", tt.contentType)

		if got, want := isJSONResponse(resp), tt.isJSON; got != want {
			t.Errorf("[%d] %q -> got %t, want %t", n, tt.contentType, got, want)
		}
	}
}
