package rails

import (
	"strings"

	"github.com/bearer/bearer/internal/report/frameworks"
)

const TypeCache frameworks.Type = "cache"
const TypeDatabase frameworks.Type = "database"
const TypeStorage frameworks.Type = "storage"

type Cache struct {
	Type string `json:"type" yaml:"type"`
}

type Database struct {
	Name    string `json:"name" yaml:"name"`
	Adapter string `json:"adapter" yaml:"adapter"`
}

type Storage struct {
	Name       string `json:"name" yaml:"name"`
	Service    string `json:"service" yaml:"service"`
	Encryption string `json:"encryption" yaml:"encryption"`
}

func (value Cache) GetTechnologyKey() string {
	switch value.Type {
	case "memory_store", "null_store":
		// Ignored cache types
		return ""
	case "file_store":
		return "39747024-c306-4a95-a0df-7e585a33a86f"
	case "mem_cache_store":
		return "42908ccc-4f0f-419e-9ba0-f25121fd15b7"
	case "redis_cache_store":
		return "62c20409-c1bf-4be9-a859-6fe6be7b11e3"
	default:
		return "unidentified_data_store"
	}
}

func (value Database) GetTechnologyKey() string {
	switch strings.ToLower(value.Adapter) {
	case "mysql2", "jdbcmysql":
		return "ffa70264-2b19-445d-a5c9-be82b64fe750"
	case "postgresql", "jdbcpostgresql":
		return "428ff7dd-22ea-4e80-8755-84c70cf460db"
	case "sqlite3", "jdbcsqlite3":
		return "aa706b3c-0f6d-4a7b-a7a5-71ee0c5b6c00"
	default:
		return "unidentified_data_store"
	}
}

func (value Storage) GetTechnologyKey() string {
	if strings.Contains(value.Name, "test") {
		// Ignored storage types
		return ""
	}

	switch value.Service {
	case "Mirror":
		// Ignored storage types
		return ""
	case "AzureStorage":
		return "f0f43ee7-7f6b-4572-aaa0-6b207146912b"
	case "Disk":
		return "39747024-c306-4a95-a0df-7e585a33a86f"
	case "GCS":
		return "3a154582-174f-4ef7-90a2-f654435c23cb"
	case "S3":
		return "4e5a3a3a-47cd-4b0e-b0a6-fa30a0a62499"
	default:
		return "unidentified_data_store"
	}
}
