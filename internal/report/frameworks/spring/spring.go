package spring

import (
	"github.com/bearer/bearer/internal/report/frameworks"
)

const TypeDatabase frameworks.Type = "database"

type DataStore struct {
	Driver string `json:"driver" yaml:"driver"`
}

func (value DataStore) GetTechnologyKey() string {
	switch value.Driver {
	case "db2", "com.ibm.db2.jcc.DB2Driver":
		return "b9bbbbb8-cb8b-4ffb-997f-e0d1e9050a96"
	case "sqlserver", "com.microsoft.sqlserver.jdbc.SQLServerDriver":
		return "e4db4505-b837-4b76-9184-c3cec3b5e522"
	case "mysql", "com.mysql.jdbc.Driver", "mariadb", "org.mariadb.jdbc.Driver":
		return "ffa70264-2b19-445d-a5c9-be82b64fe750"
	case "oracle", "oracle.jdbc.OracleDriver":
		return "80886e2a-ee2c-423d-98bc-0a3d743787b4"
	case "postgresql", "com.postgresql.jdbc.Driver":
		return "428ff7dd-22ea-4e80-8755-84c70cf460db"
	case "sqlite", "org.sqlite.JDBC":
		return "aa706b3c-0f6d-4a7b-a7a5-71ee0c5b6c00"
	default:
		return "unidentified_data_store"
	}
}
