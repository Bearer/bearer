package interfaces

import (
	"github.com/bearer/curio/pkg/report/dependencies"
)

type ClassifiedDependency struct {
	*dependencies.Dependency
	Classification Classification `json:"classification"`
}

type Classification struct {
	RecipeMatch bool
	RecipeName  string
}
