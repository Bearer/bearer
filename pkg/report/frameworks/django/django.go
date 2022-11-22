package django

import "github.com/bearer/curio/pkg/report/frameworks"

const TypeDatabase frameworks.Type = "database"

type Database struct {
	Name   string `json:"name" yaml:"name"`
	Engine string `json:"engine" yaml:"engine"`
}

func (value Database) GetTechnologyKey() string {
	switch value.Engine {
	case "sql_server.pyodbc":
		return "e4db4505-b837-4b76-9184-c3cec3b5e522"
	case "django.db.backends.mysql":
		return "ffa70264-2b19-445d-a5c9-be82b64fe750"
	case "django.db.backends.oracle":
		return "80886e2a-ee2c-423d-98bc-0a3d743787b4"
	case "django.db.backends.postgresql":
		return "428ff7dd-22ea-4e80-8755-84c70cf460db"
	case "django.db.backends.sqlite3":
		return "aa706b3c-0f6d-4a7b-a7a5-71ee0c5b6c00"
	default:
		return "unidentified_data_store"
	}
}
