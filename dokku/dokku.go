package dokku

import (
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

func New() *Dokku {
	cs := NewContainerStore()
	as := NewAppStore(cs)

	return &Dokku{
		Apps: as,
		Containers: cs,
	}
}
