package stoppable_http_server

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func BoringHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Response")
}

var once = false

func registerHandler() {
	if once == false {
		http.HandleFunc("/", BoringHandler)
		once = true
	}
}

func TestSingleThread(t *testing.T) {
	registerHandler()

	srv := StartServer()
	defer srv.Stop()

	req, err := http.Get("http://localhost:4567/")
	if err != nil {
		t.Error("Got error from server: %v\n", err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	if strings.Compare(buf.String(), "Response") != 0 || err != nil {
		t.Error("Didn't get the expected response, got \"%s\"", buf)
	}
}

func TestChangePort(t *testing.T) {
	registerHandler()
	srv := StartServer(func(config *HttpConfig) {
		config.Port = 7654
	})
	defer srv.Stop()

	req, err := http.Get("http://localhost:7654/")
	if err != nil {
		t.Error("Got error from server: %v\n", err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	if strings.Compare(buf.String(), "Response") != 0 || err != nil {
		t.Error("Didn't get the expected response, got \"%s\"", buf)
	}
}
