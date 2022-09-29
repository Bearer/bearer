package symfony

import "github.com/bearer/curio/pkg/report/frameworks"

const TypeDatabase frameworks.Type = "database"

type Database struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
}
