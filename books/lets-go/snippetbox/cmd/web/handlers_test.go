package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeeDark/go-web-labs/books/lets-go/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	// Init
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Check 1
	ping(rr, r)
	rs := rr.Result()

	// Assert 1
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// Check 2
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	// Assert 2
	assert.Equal(t, string(body), "OK")
}
