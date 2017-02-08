package stoppable_http_server

import "net/http"
import "net/url"
import "fmt"
import "sync"
import "net"

//import "log"

import "github.com/hydrogen18/stoppableListener"

type SLServer struct {
	wg     sync.WaitGroup
	sl     *stoppableListener.StoppableListener
	server *http.Server
	Url    string
}

type HttpConfig struct {
	Host string
	Port int
}

func StartServer(configFuncs ...func(*HttpConfig)) *SLServer {

	// Default configuration
	config := HttpConfig{
		Host: "127.0.0.1",
		Port: 4567,
	}

	for _, f := range configFuncs {
		f(&config)
	}

	srvIp := fmt.Sprintf("%s:%d", config.Host, config.Port)

	originalListener, err := net.Listen("tcp", srvIp)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("Starting web server at %s\n", url)

	sl, err := stoppableListener.New(originalListener)
	if err != nil {
		panic(err)
	}

	//var wg sync.WaitGroup
	srv := SLServer{server: &http.Server{},
		sl:  sl,
		wg:  sync.WaitGroup{},
		Url: fmt.Sprintf("http://%s/", srvIp)}

	srv.wg.Add(1)
	go func() {
		defer srv.wg.Done()
		srv.server.Serve(sl)
	}()

	return &srv
}

func (srv *SLServer) Stop() {
	//fmt.Printf("Stopping web server...")
	srv.sl.Stop()
	srv.wg.Wait()
	//fmt.Printf("done\n")
}

func (srv *SLServer) Wait() {
	srv.wg.Wait()
}

func (srv *SLServer) URL() url.URL {
	uri, _ := url.Parse(srv.Url)
	return *uri
}
