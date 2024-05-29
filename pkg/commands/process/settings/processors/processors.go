package processors

import (
	"embed"
	"fmt"

	"github.com/bearer/bearer/pkg/util/rego"
)

//go:embed *.rego
var fs embed.FS

func Load(name string) ([]rego.Module, error) {
	moduleText, err := fs.ReadFile(fmt.Sprintf("%s.rego", name))
	if err != nil {
		return nil, err
	}

	return []rego.Module{{
		Name:    fmt.Sprintf("bearer.%s", name),
		Content: string(moduleText),
	}}, nil
}
