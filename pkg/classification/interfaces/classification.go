package interfaces

import "github.com/bearer/curio/pkg/report/interfaces"

type ClassifiedInterface struct {
	Type           string               `json:"type"`
	Interface      interfaces.Interface `json:"interface"`
	Classification Classification       `json:"classification"`
}

type Classification struct {
}
