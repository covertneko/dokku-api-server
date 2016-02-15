// Utilities for serving HTTP over Unix sockets with Echo
package main

import (
	"net"
	"net/http"
	"os"
	"fmt"

	"github.com/labstack/echo"
)

type Echo echo.Echo

func (e *Echo) RunDomainSocket(path string, mode os.FileMode) {
	srv := &http.Server{Addr: path}
	srv.Handler = (*echo.Echo)(e)

	e.RunDomainSocketServer(srv, mode)
}

func (e *Echo) RunDomainSocketServer(srv *http.Server, mode os.FileMode) {
	// Remove old socket, if necessary/possible
	if err := os.Remove(srv.Addr); err != nil && !os.IsNotExist(err) {
		panic(fmt.Errorf("Could not remove existing unix socket at %q: %s", srv.Addr, err))
	}

	ln, err := net.Listen("unix", srv.Addr)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	if err = os.Chmod(srv.Addr, mode); err != nil {
		panic(fmt.Errorf("Could not set mode %#o for %q: %s", mode, srv.Addr, err))
	}

	if err := srv.Serve(ln); err != nil {
		ln.Close()
		panic(err)
	}
}
