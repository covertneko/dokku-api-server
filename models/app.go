package models

type App struct {
	Name			string			`json:"name"`
	Status			AppStatus		`json:"status"`
	Containers		[]*Container	`json:"-"`
	Domains			[]*Domain		`json:"-"`
}
