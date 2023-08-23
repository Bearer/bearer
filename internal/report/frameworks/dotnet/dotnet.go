package dotnet

import (
	"github.com/bearer/bearer/internal/report/frameworks"
)

const TypeDatabase frameworks.Type = "database"

type DBContext struct {
	UseDbMethodName string `json:"use_db_method_name" yaml:"use_db_method_name"`
	TypeName        string `json:"type_name" yaml:"type_name"`
}

func (value DBContext) GetTechnologyKey() string {
	switch value.UseDbMethodName {
	case "UseSqlServer":
		return "e4db4505-b837-4b76-9184-c3cec3b5e522"
	case "UseMySQL":
		return "ffa70264-2b19-445d-a5c9-be82b64fe750"
	case "UseOracle":
		return "80886e2a-ee2c-423d-98bc-0a3d743787b4"
	case "UseNpgsql":
		return "428ff7dd-22ea-4e80-8755-84c70cf460db"
	case "UseSqlite":
		return "aa706b3c-0f6d-4a7b-a7a5-71ee0c5b6c00"
	default:
		return "unidentified_data_store"
	}
}
