// Package traefik_plugin_hello return message, status code and headers with configuration and request query parameters.
package traefik_plugin_hello

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	Message    string            `json:"message,omitempty"`
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Message:    "Hello world!",
		StatusCode: 200,
		Headers:    make(map[string]string),
	}
}

// Hello a plugin.
type hello struct {
	next       http.Handler
	message    string
	statusCode int
	name       string
	headers    map[string]string
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// check here and return error in here
	return &hello{
		next:       next,
		message:    config.Message,
		statusCode: config.StatusCode,
		name:       name,
		headers:    config.Headers,
	}, nil
}

//nolint:varnamelen
func (e *hello) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	message := []byte(e.message)

	if outputQuery := req.URL.Query().Get("message"); outputQuery != "" {
		message = []byte(outputQuery)
	}

	for k, v := range e.headers {
		rw.Header().Set(k, v)
	}

	headers := req.URL.Query().Get("headers")
	if headers != "" {
		headersArray := strings.Split(req.URL.Query().Get("headers"), ",")

		for _, v := range headersArray {
			vL := strings.SplitN(v, ":", 2) //nolint:gomnd
			if len(vL) > 1 {
				rw.Header().Set(vL[0], vL[1])
			}
		}
	}

	if statusCode, err := strconv.Atoi(req.URL.Query().Get("statusCode")); err != nil {
		rw.WriteHeader(e.statusCode)
	} else {
		rw.WriteHeader(statusCode)
	}

	rw.Write(message) //nolint:gosec,errcheck
}
