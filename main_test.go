package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// test main function
var _ = func() bool {
	testing.Init()
	// set up mock command line arguments
	os.Args = append(os.Args, "parallel 3")
	os.Args = append(os.Args, "google.com")
	os.Args = append(os.Args, "adjust.com")

	return true
}()

func TestGetMD5Hash(t *testing.T) {
	// test MD5 hash of "test"
	expectedHash := "098f6bcd4621d373cade4e832627b4f6" // MD5 hash of "test"

	// get MD5 hash of "test"
	hash := GetMD5Hash("test")

	// check hash
	if hash != expectedHash {
		t.Errorf("Expected MD5 hash: %s, but got: %s", expectedHash, hash)
	}
}

func TestParallelGet(t *testing.T) {

	// set up mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}))
	defer mockServer.Close()
	for i := range urls {
		urls[i] = strings.Split(mockServer.URL, "http://")[1]
	}

	// run parallelGet
	results := parallelGet(urls, 2)

	// check results
	if len(urls) != len(results) {
		t.Errorf("Incorrect number of results: %d, but got: %d", len(urls), len(results))
	}

	// check results
	for _, result := range results {
		if mockServer.URL != result.url {
			t.Errorf("URL mismatch")
		}
		if "098f6bcd4621d373cade4e832627b4f6" != result.hashResponse {
			t.Errorf("Hash response mismatch")
		}
	}
}
