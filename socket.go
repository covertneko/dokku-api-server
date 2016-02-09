// Utilities for serving HTTP over Unix sockets
package main

import (
	"net"
	"net/http"
	"os"
	"fmt"
)

type Server http.Server

// Serve HTTP requests on a Unix socket with a given file mode
// Pretty much copied from the implementation in
// github.com/valyala/fasthttp/server.go
func ListenAndServeUnix(addr string, mode os.FileMode, handler http.Handler) error {
	srv := &Server {
		Addr: addr,
		Handler: handler,
	}

	return srv.ListenAndServeUnix(mode)
}

func (srv *Server) ListenAndServeUnix(mode os.FileMode) error {
	// Remove old socket, if necessary/possible
	if err := os.Remove(srv.Addr); err != nil && !os.IsNotExist(err) {
		fmt.Errorf("Could not remove existing unix socket at %q: %s", srv.Addr, err)
	}

	l, err := net.Listen("unix", srv.Addr)
	if err != nil {
		return err
	}

	if err = os.Chmod(srv.Addr, mode); err != nil {
		return fmt.Errorf("Could not set mode %#o for %q: %s", mode, srv.Addr, err)
	}

	return (*http.Server)(srv).Serve(l)
}
