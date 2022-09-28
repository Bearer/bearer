package dotnet

import (
	"github.com/bearer/curio/pkg/report/frameworks"
)

const TypeDatabase frameworks.Type = "database"

type DBContext struct {
	UseDbMethodName string `json:"use_db_method_name"`
	TypeName        string `json:"type_name"`
}
