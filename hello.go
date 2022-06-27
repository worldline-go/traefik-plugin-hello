package traefik_plugin_hello

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

var defaultMsg = "Hello world!"

// Config the plugin configuration.
type Config struct {
	Message    string            `json:"message,omitempty"`
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Message:    defaultMsg,
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

func prepareOutput(req *http.Request, msg string) []byte {
	result := make(map[string]interface{})

	result["msg"] = msg
	result["header"] = req.Header
	result["URL"] = req.URL.String()

	body := ""

	if req.Body != nil {
		defer req.Body.Close()

		buf := bytes.Buffer{}
		buf.ReadFrom(req.Body) //nolint:errcheck
		body = buf.String()
	}

	result["body"] = body
	result["form"] = req.Form.Encode()

	result["host"] = req.Host
	result["remoteAddr"] = req.RemoteAddr
	result["requestURI"] = req.RequestURI
	result["TLS"] = req.TLS
	result["trailer"] = req.Trailer
	result["transferEncoding"] = req.TransferEncoding

	bf := bytes.NewBufferString("")
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(result) //nolint:errcheck

	return bf.Bytes()
}

func (e *hello) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	message := []byte(e.message)

	if outputQuery := req.URL.Query().Get("message"); outputQuery != "" {
		message = []byte(outputQuery)
	}

	details, _ := strconv.ParseBool(req.URL.Query().Get("details"))

	if details {
		message = prepareOutput(req, string(message))

		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
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

	rw.Write(message) //nolint:errcheck
}
