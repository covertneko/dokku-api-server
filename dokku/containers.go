package dokku

import (
	"fmt"
	"log"
	"strings"
	"sync"

	m "github.com/nikelmwann/dokku-api-server/models"
	"github.com/patrickmn/go-cache"
)

type ContainerStore struct {
	containers *cache.Cache
}

func NewContainerStore() *ContainerStore {
	return &ContainerStore{containers: newCache()}
}

func (s *ContainerStore) Invalidate(id string) {
	_, cached := s.containers.Get(id)
	if !cached {
		// Nothing to invalidate
		return
	}

	log.Println("Invalidating cache entry for container ", id)
	s.containers.Delete(id)
}

// Find a container by inspecting dokku command output
func (s *ContainerStore) lookup(id string) (container *m.Container, err error) {
	// Gather output of `dokku ls` to read container list
	out, err := Exec("ls")
	if err != nil {
		err = fmt.Errorf("Error requesting container %q: %s", id, err)
		return
	}

	// TODO: find a better way to do this (maybe by inspecting the app folder?)
	for _, line := range out {
		if !strings.Contains(line, id) {
			// Line does not contain requested container; skip it
			continue
		}

		var cerr error
		container, cerr = m.ParseContainer(line)
		if cerr == nil {
			log.Printf(
				"Found container %q; type %q; %s",
				container.ID, container.Type, container.Status)
			return
		}
	}

	// If control reaches this stage, the container was not found
	err = fmt.Errorf("Could not find container %q", id)
	return
}

// Find container
func (s *ContainerStore) Find(id string) (container *m.Container, err error) {
	log.Printf("Finding container %q...", id)

	item, cached := s.containers.Get(id)
	if !cached {
		log.Printf("Container %q not found in cache", id)
		container, err = s.lookup(id)
		if err != nil {
			err = fmt.Errorf("Error requesting container %q: %s", id, err)
		}

		// Store pointer to container in cache
		s.containers.Set(id, container, cache.DefaultExpiration)
	} else {
		container, cached = item.(*m.Container)
		if !cached {
			err = fmt.Errorf("Type error when retrieving container %q from cache", id)
		}
	}

	return
}

// Find all containers
func (s *ContainerStore) FindAll() (containers []*m.Container, err error) {
	output, err := Exec("ls")
	if err != nil {
		err = fmt.Errorf("Error requesting containers list: %s", err)
		return
	}

	var wg sync.WaitGroup
	// Skip first line
	// Remaining output is the list of containers; one container per line
	for _, line := range output[1:] {
		wg.Add(1)
		go func() {
			defer wg.Done()

			c, cerr := s.Find(line)
			if cerr != nil {
				panic(fmt.Errorf("Error requesting containers list: %s", cerr))
			}

			containers = append(containers, c)
		}()
	}

	wg.Wait()

	return
}
