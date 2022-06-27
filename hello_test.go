package traefik_plugin_hello_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	traefik_plugin_hello "github.com/worldline-go/traefik-plugin-hello"
)

func TestHello(t *testing.T) {
	cfg := traefik_plugin_hello.CreateConfig()
	cfg.Headers["Content-type"] = "application/json"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := traefik_plugin_hello.New(ctx, next, cfg, "hello-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, recorder, "Content-Type", []string{"application/json"})
	if recorder.Code != 200 {
		t.Errorf("Expected status code 200, got %d", recorder.Code)
	}
}

func assertHeader(t *testing.T, rw *httptest.ResponseRecorder, key string, expected []string) {
	t.Helper()

	if !reflect.DeepEqual(expected, rw.Header()[key]) {
		t.Errorf("invalid header value: %v", rw.Header()[key])
	}
}
