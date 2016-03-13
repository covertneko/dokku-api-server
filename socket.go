// Utilities for serving HTTP over Unix sockets with Echo
package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

type Echo echo.Echo

func (e *Echo) RunDomainSocket(path string, mode os.FileMode) error {
	srv := &http.Server{Addr: path}
	srv.Handler = (*echo.Echo)(e)

	return e.RunDomainSocketServer(srv, mode)
}

func (e *Echo) RunDomainSocketServer(srv *http.Server, mode os.FileMode) error {
	ln, err := ListenSocket(srv.Addr, mode)
	if err != nil {
		return err
	}
	defer ln.Close()
	defer os.Remove(srv.Addr)

	if err = srv.Serve(ln); err != nil {
		return err
	}

	return nil
}

func ListenSocket(addr string, mode os.FileMode) (net.Listener, error) {
	// Remove old socket, if necessary/possible
	if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf(
			"Could not remove existing unix socket at %q: %s", addr, err)
	}

	ln, err := net.Listen("unix", addr)
	if err != nil {
		return nil, err
	}

	if err = os.Chmod(addr, mode); err != nil {
		return nil, fmt.Errorf(
			"Could not set mode %#o for %q: %s", mode, addr, err)
	}

	return ln, nil
}
