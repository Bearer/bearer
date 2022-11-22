package django

import "github.com/bearer/curio/pkg/report/frameworks"

const TypeDatabase frameworks.Type = "database"

type Database struct {
	Name   string `json:"name" yaml:"name"`
	Engine string `json:"engine" yaml:"engine"`
}
