package traefik_plugin_hello_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	hello "github.com/worldline-go/traefik-plugin-hello"
)

func TestHello(t *testing.T) {
	cfg := hello.CreateConfig()
	cfg.Headers["Content-type"] = "application/json"

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := hello.New(ctx, next, cfg, "hello-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	fmt.Println(recorder.Header())

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
