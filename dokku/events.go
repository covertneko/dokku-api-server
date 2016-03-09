package dokku

import (
	"fmt"
	"strings"
)

//go:generate stringer -type=EventType

type EventType int

// TODO: Add all dokku events
const (
	PreDeploy EventType = iota
	PostDeploy
	CheckDeploy
	PreDelete
	PostDelete
	RecieveApp
	PostDomainsUpdate
)

type Event struct {
	AppName string
	Type    EventType
}

func ParseEvent(line string) (*Event, error) {
	// Split output in half - important data is on the right side of the token
	cols := strings.Split(line, "INVOKED: ")
	ncols := len(cols)

	if ncols != 2 {
		return nil, fmt.Errorf("Invalid entry in Dokku event log: %q", line)
	}

	// Event name and app name are the last column
	// Displayed in the form:
	// "event_name( app_name container_id container_type additional_info )"
	// Right now all we care about is the event and app names

	// Get event name
	cols = strings.Split(cols[1], "( ")
	name := cols[0]

	// Get app name
	cols = strings.Split(cols[1], " ")
	app := cols[0]

	// Get event type from event name
	var t EventType
	switch name {
	case "pre-deploy":
		t = PreDeploy
	case "post-deploy":
		t = PostDeploy
	case "check-deploy":
		t = CheckDeploy
	case "pre-delete":
		t = PreDelete
	case "post-delete":
		t = PostDelete
	case "recieve-app":
		t = RecieveApp
	case "post-domains-update":
		t = PostDomainsUpdate
	default:
		return nil, fmt.Errorf("Unsupported event %q", name)
	}

	return &Event{
		AppName: app,
		Type:    t,
	}, nil
}
