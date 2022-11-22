package rails

import (
	"strings"

	"github.com/bearer/curio/pkg/report/frameworks"
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
		return "Disk"
	case "mem_cache_store":
		return "Memcached"
	case "redis_cache_store":
		return "Redis"
	default:
		return "unidentified_data_store"
	}
}

func (value Database) GetTechnologyKey() string {
	switch strings.ToLower(value.Adapter) {
	case "mysql2", "jdbcmysql":
		return "MySQL"
	case "postgresql", "jdbcpostgresql":
		return "PostgreSQL"
	case "sqlite3", "jdbcsqlite3":
		return "SQLite"
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
		return "Azure Storage"
	case "Disk":
		return "Disk"
	case "GCS":
		return "Google Cloud Storage"
	case "S3":
		return "AWS S3"
	default:
		return "unidentified_data_store"
	}
}
