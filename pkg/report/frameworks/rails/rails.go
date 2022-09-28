package rails

import "github.com/bearer/curio/pkg/report/frameworks"

const TypeCache frameworks.Type = "cache"
const TypeDatabase frameworks.Type = "database"
const TypeStorage frameworks.Type = "storage"

type Cache struct {
	Type string `json:"type"`
}

type Database struct {
	Name    string `json:"name"`
	Adapter string `json:"adapter"`
}

type Storage struct {
	Name       string `json:"name"`
	Service    string `json:"service"`
	Encryption string `json:"encryption"`
}
