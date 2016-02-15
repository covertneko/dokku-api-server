package main

import (
	"os"
	"path"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/nikelmwann/dokku-api/resources"
)

func main() {
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Get("/", Index)
	e.Get("/apps", resources.HandlerFor(resources.Apps{}))
	e.Get("/apps/:name", resources.HandlerFor(resources.App{}))

	// Get socket path from environment or default
	var sockpath string
	if p := os.Getenv("DOKKU_API_SOCKET"); p != "" {
		sockpath = p
	} else {
		sockpath = "/tmp/dokku-api/api.sock"
	}

	// Ensure socket directory exists
	if _, err := os.Stat(sockpath); os.IsNotExist(err) {
		os.MkdirAll(path.Dir(sockpath), 0777)
	}

	(*Echo)(e).RunDomainSocket(sockpath, os.FileMode(0666))
}
