// models/Role.go
package models

type Role int

const (
	RoleReader Role = iota
	RoleAdmin
)
