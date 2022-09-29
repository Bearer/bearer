package schema

import (
	"github.com/bearer/curio/pkg/report/schema"
)

type ClassifiedSchema struct {
	Type           string         `json:"type"`
	Schema         schema.Schema  `json:"schema"`
	Classification Classification `json:"classification"`
}

type Classification struct {
}
