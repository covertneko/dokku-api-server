package models

//go:generate jsonenums -type=AppStatus

type AppStatus int

const (
	Running AppStatus = iota
	Stopped
	Undeployed
)
