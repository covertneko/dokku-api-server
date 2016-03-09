package models

import (
	"fmt"
	"strings"
)

type Container struct {
	ID string			`json:"id"`
	App string			`json:"app"`
	Status AppStatus	`json:"status"`
	Type string			`json:"type"`
}

// Implements dokku.Cachable
func (c *Container) GetID() string {
	return c.ID
}

// Attempt to parse a container from a line of text
func ParseContainer(line string) (c *Container, err error) {
	// `dokku ls` displays containers in four columns like so:
	// appname    type    id    status
	cols := strings.Fields(line)

	// If this line does not describe a container, skip it
	if len(cols) < 4 {
		err = fmt.Errorf("Line does not describe a container.")
		return
	}

	c = new(Container)
	c.App = cols[0]
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

	return
}
