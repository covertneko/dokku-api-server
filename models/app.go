package models

import (
	"strings"

	"github.com/nikelmwann/dokku-api/dokku"
)

type App struct {
	Name string `json:"name"`
	Running bool `json:"running"`
}
type Apps []App

func GetApp(name string) (app App, err error) {
	out, err := dokku.Exec("ps", name)

	running := false
	// Check output of first line for "running".
	// TODO: is there a better way to detect this?
	// For now, if the container is running, the app will be considered to be
	// running.
	if strings.Contains(<-out, "running") {
		running = true
	}

	app = App{Name: name, Running: running}
	return
}
