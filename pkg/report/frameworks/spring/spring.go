package spring

import "github.com/bearer/curio/pkg/report/frameworks"

const TypeDatabase frameworks.Type = "database"

type DataStore struct {
	Driver string `json:"driver" yaml:"driver"`
}
