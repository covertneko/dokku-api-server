package models

import (
	"fmt"
	"strings"

	"github.com/nikelmwann/dokku-api/dokku"
)

type App struct {
	Name			string			`json:"name"`
	Status			AppStatus		`json:"status"`
	Containers		[]*Container	`json:"containers"`
	Domains			[]*Domain		`json:"domains"`
}

// Attempt to parse a container for a given app from a line
func parseContainer(name string, line string) (c Container, err error) {
	// `dokku ls` displays containers in four columns like so:
	// appname    type    id    status
	cols := strings.Fields(line)

	// If this line does not describe a container, skip it
	if len(cols) < 4 {
		err = fmt.Errorf("Line does not describe a container.")
		return
	}

	// Construct a Container from the line if it belongs to the app
	if cols[0] == name {
		c.Type = cols[1]
		c.ID = cols[2]

		switch cols[3] {
		case "running":
			c.Status = Running
		case "stopped":
			c.Status = Stopped
		default:
			c.Status = Undeployed
		}
	} else {
		err = fmt.Errorf("Line does not contain a container for app %q", name)
	}

	return
}

// Find app and all containers
func FindApp(name string) (app App, err error) {
	app.Name = name
	// Set initial app status to undeployed
	// Status will be updated based on container statuses
	app.Status = Undeployed

	out, err := dokku.Exec("ls")
	if err != nil {
		err = fmt.Errorf("Error requesting app %q: %s", name, err)
		return
	}

	// Find all containers for app
	for _, line := range out {
		fmt.Println(line)

		c, cerr := parseContainer(name, line)
		if cerr == nil {
			app.Containers = append(app.Containers, &c)

			if c.Status == Running {
				app.Status = Running
			}

			if c.Status == Stopped && app.Status == Undeployed {
				app.Status = Stopped
			}
		}
	}

	return
}

func FindApps() (apps []App, err error) {
	output, err := dokku.Exec("apps")
	if err != nil {
		err = fmt.Errorf("Error requesting apps list: %s", err)
		return
	}

	// Remaining output is the list of apps; one app per line
	for _, line := range output[1:] {
		var app App
		app, err = FindApp(line)
		if err != nil {
			err = fmt.Errorf("Error requesting apps list: %s", err)
			return
		}

		apps = append(apps, app)
	}

	return
}
