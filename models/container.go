package models

type Container struct {
	ID string			`json:"-"`
	Status AppStatus	`json:"status"`
	Type string			`json:"type"`
}
