package dokku

import (
	"fmt"
	"log"
	"time"
	"github.com/patrickmn/go-cache"
)

const (
	DOKKU_EVENT_LOG = "/var/log/dokku/events.log"
	CACHE_EXPIRATION_INTERVAL = 5*time.Minute
	CACHE_CLEANUP_INTERVAL = 30*time.Second
)

func newCache() *cache.Cache {
	return cache.New(CACHE_EXPIRATION_INTERVAL, CACHE_CLEANUP_INTERVAL)
}

type Cachable interface {
	GetID() string
}

type Store interface {
	Invalidate(id string)
	Find(id string) (*Cachable, error)
	FindAll() ([]*Cachable, error)
}

type Dokku struct {
	Apps       *AppStore
	Containers *ContainerStore
}

// Monitor dokku events and invalidate caches when necessary
func (d *Dokku) processEvents() {
	// Make sure events are enabled
	_, err := Exec("events:on")
	if err != nil {
		panic(fmt.Errorf("Could not open Dokku events log: %s", err))
	}

	// Monitor events by following the dokku event log.
	// -n0 is necessary to prevent old events from being processed
	output, err := followCmd("tail", "-fn0", DOKKU_EVENT_LOG)
	if err != nil {
		panic(fmt.Errorf("Could not open Dokku events log: %s", err))
	}

	// Process events
	for ln := range output.Lines {
		e, err := ParseEvent(ln)
		if err != nil {
			// Unsupported event; skip
			log.Println(err)
			continue
		}

		log.Printf("Got event %s for app %q\n", e.Type, e.AppName)

		// Invalidate app
		d.Apps.Invalidate(e.AppName)
	}
}

func New() *Dokku {
	cs := NewContainerStore()
	as := NewAppStore(cs)

	d := &Dokku{
		Apps: as,
		Containers: cs,
	}

	go d.processEvents()

	return d
}
