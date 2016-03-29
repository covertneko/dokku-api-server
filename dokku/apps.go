package dokku

import (
	"fmt"
	"log"

	m "github.com/nikelmwann/dokku-api-server/models"
	"github.com/patrickmn/go-cache"
)

type AppStore struct {
	apps           *cache.Cache
	containerStore *ContainerStore
}

func NewAppStore(containerStore *ContainerStore) *AppStore {
	return &AppStore{apps: newCache(), containerStore: containerStore}
}

func (s *AppStore) Invalidate(id string) {
	_app, cached := s.apps.Get(id)
	if !cached {
		// Nothing to invalidate
		return
	}
	app := _app.(*m.App)

	log.Println("Invalidating cache entry for app ", id)

	// Invalidate all containers for app first
	for _, c := range app.Containers {
		s.containerStore.Invalidate(c)
	}

	// Then invalidate app
	s.apps.Delete(id)
}

// Find an app by inspecting dokku command output
func (s *AppStore) lookup(name string) (app *m.App, err error) {
	// First ensure the app exists
	out, err := Exec("apps")
	if err != nil {
		err = fmt.Errorf("Error requesting app %q: %s", name, err)
		return
	}

	exists := false
	for _, line := range out {
		if line == name {
			exists = true
		}
	}

	if !exists {
		// App does not exist; return nil
		return
	}

	// App exists; start gathering additional info
	app = new(m.App)
	app.Name = name
	// Set initial app status to undeployed
	// Status will be updated based on container statuses
	app.Status = m.Undeployed

	// Gather output of `dokku ls` to read container list
	out, err = Exec("ls")
	if err != nil {
		err = fmt.Errorf("Error requesting app %q: %s", name, err)
		return
	}

	// Find all containers for app
	// TODO: find a better way to do this (maybe by inspecting the app folder?)
	for _, line := range out {
		c, cerr := m.ParseContainer(line)
		if cerr == nil && c.App == name {
			log.Printf("Found container %q; type %q; %s", c.ID, c.Type, c.Status)
			app.Containers = append(app.Containers, c.GetID())

			if c.Status == m.Running {
				app.Status = m.Running
			}

			if c.Status == m.Stopped && app.Status == m.Undeployed {
				app.Status = m.Stopped
			}
		}
	}

	// Find domains for app
	out, err = Exec("domains", name)
	if err != nil {
		err = fmt.Errorf("Error requesting domains for app %q: %s", name, err)
		return
	}

	for _, line := range out {
		d, derr := m.NewDomain(line)
		if derr == nil {
			log.Printf("Found domain %q", d.Name)
			app.Domains = append(app.Domains, d.Name)
		}
	}

	return
}

// Find app and all containers
func (s *AppStore) Find(id string) (app *m.App, err error) {
	log.Printf("Finding app %q...", id)

	item, cached := s.apps.Get(id)
	if !cached {
		log.Printf("App %q not found in cache", id)
		app, err = s.lookup(id)
		if err != nil {
			err = fmt.Errorf("Error requesting app %q: %s", id, err)
		}

		// Store pointer to app in cache
		s.apps.Set(id, app, cache.DefaultExpiration)
	} else {
		app, cached = item.(*m.App)
		if !cached {
			err = fmt.Errorf("Type error when retrieving app %q from cache", id)
		}
	}

	return
}

// Find all apps
func (s *AppStore) FindAll() (apps []*m.App, err error) {
	output, err := Exec("apps")
	if err != nil {
		err = fmt.Errorf("Error requesting apps list: %s", err)
		return
	}

	// Skip first line
	// Remaining output is the list of apps; one app per line
	for _, line := range output[1:] {
		var app *m.App
		app, err = s.Find(line)
		if err != nil {
			err = fmt.Errorf("Error requesting apps list: %s", err)
			return
		}

		apps = append(apps, app)
	}

	return
}
