package models

//go:generate jsonenums -type=AppStatus
//go:generate stringer -type=AppStatus

type AppStatus int

const (
	Running AppStatus = iota
	Stopped
	Undeployed
)
