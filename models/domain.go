package models

import (
	"fmt"
	"regexp"
)

type Domain struct {
	Name string `json:"name"`
}

// Attempt to parse a domain from a line
func NewDomain(line string) (d *Domain, err error) {
	// Mmmmm domain name regex
	re := regexp.MustCompile(`^([a-z0-9][a-z0-9-]*\.)*[a-z0-9][a-z0-9-]*\.?$`)

	// If this line does not describe a valid domain, skip it
	if !re.MatchString(line) {
		err = fmt.Errorf("Line does not describe a domain.")
		return
	}

	d = new(Domain)
	d.Name = line
	return
}
