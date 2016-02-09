// Utilities for using HTTP over Unix sockets
// Mostly copied from this dead PR:
// https://github.com/apatil/napping-unixsocket/blob/master/unix_socket.go

package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"path"
	"os"
	"fmt"
)

type Server http.Server

// Serve HTTP requests on a Unix socket with a given file mode
// Pretty much copied from the implementation in
// github.com/valyala/fasthttp/server.go
func ListenAndServeUNIX(addr string, mode os.FileMode, handler http.Handler) error {
	srv := &Server {
		Addr: addr,
		Handler: handler,
	}

	return srv.ListenAndServeUNIX(mode)
}

func (srv *Server) ListenAndServeUNIX(mode os.FileMode) error {
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

// Transport for HTTP requests over sockets
type SocketTransport struct { path string }

// Roundtripper for unix socket requests
func (t SocketTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	dial, err := net.Dial("unix", t.path)
	if err != nil {
		return nil, err
	}

	conn := httputil.NewClientConn(dial, nil)
	defer conn.Close()
	return conn.Do(req)
}

// Helper to test if a path identifies a unix socket
func isSocket(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return fi.Mode() & os.ModeType == os.ModeSocket
}

// Split a path into a socket path and request path.
// Returns an error if the path does not identify a socket.
func GetSocketRequest(rawPath string) (string, string, error) {
	p := rawPath
	// Ensure path is absolute
	if p[0] != '/' {
		p = "/" + p
	}

	req := ""
	req_ := ""
	for p != "" {
		// Remove trailing slash from path, if any
		if l := len(p) - 1; l >= 0 && p[l] == '/' {
			p = p[:l]
		}

		if isSocket(p) {
			return p, "/" + req, nil
		}

		// Path is not a socket. Prepend path node to request, set p to
		// remaining path.
		p, req_ = path.Split(p)
		req = path.Join(req_, req)
	}

	return "", "", fmt.Errorf("%q does not identify a socket", rawPath)
}
