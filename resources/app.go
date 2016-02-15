package resources

import (
	"fmt"
	"strings"
	"net/http"

	"github.com/labstack/echo"
	"github.com/nikelmwann/dokku-api/dokku"
	"github.com/nikelmwann/dokku-api/models"
)

type Apps struct {
	PostNotSupported
	PutNotSupported
	DeleteNotSupported
}

// Find all apps
func (r Apps) Find() (apps []models.App, err error) {
	output, err := dokku.Exec("apps")
	if err != nil {
		err = fmt.Errorf("Error requesting apps list: ", err)
		return
	}

	// Skip first line of output which is simply "=====> My Apps"
	<-output

	// Remaining output is the list of apps; one app per line
	for line := range output {
		var app models.App
		ar := App{}
		app, err = ar.Find(line)
		if err != nil {
			err = fmt.Errorf("Error requesting apps list: ", err)
			return
		}

		apps = append(apps, app)
	}

	return
}

func (r Apps) Get(c *echo.Context) error {
	apps, err := r.Find()
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

// Find app and all containers
func (App) Find(name string) (app models.App, err error) {
	out, err := dokku.Exec("ls")

	// Find all containers for app
	for line := range out {
		// `dokku ls` displays containers in four columns like so:
		// appname    type    id    status
		cols := strings.Fields(line)

		// If this line does not describe a container, skip it
		if len(cols) < 4 {
			continue
		}

		// Construct a Container from the line if it belongs to the app
		if cols[0] == name {
			var c models.Container

			switch cols[3] {
			case "running":
				c.Status = models.Running
			case "stopped":
				c.Status = models.Stopped
			default:
				c.Status = models.Undeployed
			}

			if c.Status == models.Undeployed {
			}
		}
	}

	app.Name = name
	return
}

func (r App) Get(c *echo.Context) error {
	name := c.P(0)

	app, err := r.Find(name)
	if err != nil {
		return err
	}

	return c.JSONIndent(http.StatusOK, app, "", "  ")
}
