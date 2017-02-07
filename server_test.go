package stoppable_http_server

import (
  "testing"
  "net/http"
  "io"
  "bytes"
  "strings"
)

func BoringHandler( w http.ResponseWriter, r *http.Request) {
  io.WriteString( w, "Response" )
}

func TestSingleThread( t *testing.T ) {
  http.HandleFunc("/", BoringHandler )
  srv := HttpServer()
  defer srv.Stop()

  req,err := http.Get("http://localhost:4567/")
  if err != nil {
    t.Error("Got error from server: %v\n", err )
  }
  buf := new(bytes.Buffer)
  buf.ReadFrom( req.Body )
  if strings.Compare( buf.String(), "Response") != 0 || err != nil {
    t.Error("Didn't get the expected response, got \"%s\"", buf )
  }

}
