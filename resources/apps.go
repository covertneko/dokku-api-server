package resources

import (
	"net/http"
	//"strings"

	"github.com/labstack/echo"
	"github.com/nikelmwann/dokku-api-server/dokku"
)

type Apps struct {
	PostNotSupported
	PutNotSupported
	DeleteNotSupported
}

// Find all apps
func (r Apps) Get(c *echo.Context, s *dokku.Dokku) error {
	// fields := strings.Split(c.Param("fields"), ",")

	apps, err := s.Apps.FindAll()
	if err != nil {
		return err
	}

	return c.JSONIndent(http.StatusOK, apps, "", "  ")
}

type App struct {
	PostNotSupported
	PutNotSupported
	DeleteNotSupported
}

// Find one app
func (r App) Get(c *echo.Context, s *dokku.Dokku) error {
	name := c.P(0)

	app, err := s.Apps.Find(name)
	if err != nil {
		return err
	}

	if app == nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSONIndent(http.StatusOK, app, "", "  ")
}
