package models

type Container struct {
	ID string			`json:"id"`
	Status AppStatus	`json:"status"`
	Type string			`json:"type"`
}
