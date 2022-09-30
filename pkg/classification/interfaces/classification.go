package interfaces

import "github.com/bearer/curio/pkg/report/interfaces"

type ClassifiedInterface struct {
	*interfaces.Interface
	Classification Classification `json:"classification"`
}

type Classification struct {
}
