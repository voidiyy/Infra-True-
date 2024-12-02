package config

import (
	"fmt"
)

const (
	DB         = "DB"
	Cache      = "Cache"
	ServerMain = "Server_main"
	ServerAdd  = "Server_add"
	LB         = "LB"
	Other      = "Other"
)

// ContainerTypes is a list of valid container types
type ContainerTypes struct {
	types map[string]bool
}

// NewContainerTypes returns a new instance of ContainerTypes
func NewContainerTypes() *ContainerTypes {
	return &ContainerTypes{
		types: map[string]bool{
			"DB":          true,
			"Cache":       true,
			"Server_main": true,
			"Server_add":  true,
			"LB":          true,
			"Other":       true,
		},
	}
}

// AddType adds a new container type
func (ct *ContainerTypes) AddType(tp string) error {
	if tp == "" || ct.types[tp] {
		return fmt.Errorf("invalid or duplicate container type: %s", tp)
	}
	ct.types[tp] = true
	return nil
}

// IsValid checks if a container type is valid
func (ct *ContainerTypes) IsValid(tp string) bool {
	return ct.types[tp]
}
