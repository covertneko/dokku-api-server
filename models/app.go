package models

type App struct {
	Name			string		`json:"name"`
	Status			AppStatus	`json:"status"`
	Containers		[]string	`json:"containers,omitempty"`
	Domains			[]string	`json:"domains,omitempty"`
	Configuration	[]string	`json:"configuration,omitempty"`
}

// Implements dokku.Cachable
func (a * App) GetID() string {
	return a.Name
}
