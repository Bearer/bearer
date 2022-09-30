package schema

import (
	"github.com/bearer/curio/pkg/report/schema"
)

type ClassifiedSchema struct {
	*schema.Schema
	Classification Classification `json:"classification"`
}

type Classification struct {
}
