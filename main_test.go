package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func makeRequest(method, path string) *http.Response {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	rbay(w, req)
	return w.Result()
}

func suppressLogOutput() func() {
	defaultLogger := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	return func() {
		slog.SetDefault(defaultLogger)
	}
}
func TestRootReturns200(t *testing.T) {
	restoreLog := suppressLogOutput()
	defer restoreLog()
	res := makeRequest(http.MethodGet, "/")
	expected := 200
	if res.StatusCode != expected {
		t.Fatal(fmt.Sprintf("Expected status of %d got %d", expected, res.StatusCode))
	}
}

func TestRootReturnsGetMethod(t *testing.T) {
	res := makeRequest(http.MethodGet, "/")
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := "Method: GET"
	if !strings.Contains(string(data), expected) {
		t.Fatal(fmt.Sprintf("Expected field '%s' but it was not found in the response", expected))
	}
}

func TestRootReturnsPostMethod(t *testing.T) {
	res := makeRequest(http.MethodPost, "/")
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := "Method: POST"
	if !strings.Contains(string(data), expected) {
		t.Fatal(fmt.Sprintf("Expected field '%s' but it was not found in the response", expected))
	}
}

func TestRootReturnsPath(t *testing.T) {
	res := makeRequest(http.MethodPost, "/foo")
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := "URL: /foo"
	if !strings.Contains(string(data), expected) {
		t.Fatal(fmt.Sprintf("Expected field '%s' but it was not found in the response", expected))
	}
}

func TestHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("foo", "bar")
	w := httptest.NewRecorder()
	rbay(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	expected := "Headers:"
	if !strings.Contains(string(data), expected) {
		t.Fatal(fmt.Sprintf("Expected field '%s' but it was not found in the response", expected))
	}
	expected = "Foo: [bar]"
	if !strings.Contains(string(data), expected) {
		t.Fatal(fmt.Sprintf("Expected field '%s' but it was not found in the response", expected))
	}
}

func TestDefaultMessage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/message", nil)
	w := httptest.NewRecorder()
	os.Unsetenv("MESSAGE")
	message(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != DEFAULT_MESSAGE {
		t.Errorf("Expected message '%s' but got '%s'\n", DEFAULT_MESSAGE, data)
	}
}

func TestCustomMessage(t *testing.T) {
	customMessage := "foo bar"
	req := httptest.NewRequest(http.MethodGet, "/message", nil)
	w := httptest.NewRecorder()
	os.Setenv("MESSAGE", customMessage)
	message(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if string(data) != customMessage {
		t.Errorf("Expected message '%s' but got '%s'\n", customMessage, data)
	}
}

func TestMessageWrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/message", nil)
	w := httptest.NewRecorder()
	message(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code '%d' but got '%d'\n", http.StatusMethodNotAllowed, res.StatusCode)
	}
}
