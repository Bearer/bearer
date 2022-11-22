package beego

import (
	"github.com/bearer/curio/pkg/report/frameworks"
)

const TypeDatabase frameworks.Type = "database"

type Database struct {
	Name         string `json:"name" yaml:"name"`
	DriverName   string `json:"driver_name" yaml:"driver_name"`
	Package      string `json:"package" yaml:"package"`
	TypeConstant string `json:"type_constant" yaml:"type_constant"`
}
