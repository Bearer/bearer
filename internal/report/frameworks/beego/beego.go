package beego

import (
	"github.com/bearer/bearer/internal/report/frameworks"
)

const TypeDatabase frameworks.Type = "database"

type Database struct {
	Name         string `json:"name" yaml:"name"`
	DriverName   string `json:"driver_name" yaml:"driver_name"`
	Package      string `json:"package" yaml:"package"`
	TypeConstant string `json:"type_constant" yaml:"type_constant"`
}

func (value Database) GetTechnologyKey() string {
	if value.Package != "" {
		return technologyForDriverLib(value.Package, value.TypeConstant)
	}

	return technologyForDriverName(value.DriverName)
}

func technologyForDriverLib(packageName string, typeConstant string) string {
	if packageName != "github.com/beego/beego/v2/client/orm" {
		return "unidentified_data_store"
	}

	switch typeConstant {
	case "DRMySQL", "DR_MySQL":
		return "ffa70264-2b19-445d-a5c9-be82b64fe750"
	case "DRPostgres", "DR_Postgres":
		return "428ff7dd-22ea-4e80-8755-84c70cf460db"
	case "DRSqlite", "DR_Sqlite":
		return "aa706b3c-0f6d-4a7b-a7a5-71ee0c5b6c00"
	default:
		return "unidentified_data_store"
	}
}

func technologyForDriverName(driverName string) string {
	switch driverName {
	case "mysql":
		return "ffa70264-2b19-445d-a5c9-be82b64fe750"
	case "postgres":
		return "428ff7dd-22ea-4e80-8755-84c70cf460db"
	case "sqlite3":
		return "aa706b3c-0f6d-4a7b-a7a5-71ee0c5b6c00"
	default:
		return "unidentified_data_store"
	}
}
