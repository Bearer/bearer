package beego

import (
	"github.com/bearer/curio/pkg/report/frameworks"
)

const TypeDatabase frameworks.Type = "database"

type Database struct {
	Name         string `json:"name"`
	DriverName   string `json:"driver_name"`
	Package      string `json:"package"`
	TypeConstant string `json:"type_constant"`
}
