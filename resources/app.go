package resources

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/nikelmwann/dokku-api/models"
)

type Apps struct {
	PostNotSupported
	PutNotSupported
	DeleteNotSupported
}

// Find all apps
func (r Apps) Get(c *echo.Context) error {
	apps, err := models.FindApps()
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

func (r App) Get(c *echo.Context) error {
	name := c.P(0)

	app, err := models.FindApp(name)
	if err != nil {
		return err
	}

	return c.JSONIndent(http.StatusOK, app, "", "  ")
}
