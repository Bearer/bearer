package language

import (
	"fmt"

	"github.com/bearer/curio/new/language/base"
	"github.com/bearer/curio/new/language/implementation"
	"github.com/bearer/curio/new/language/implementation/ruby"
	"github.com/bearer/curio/new/language/types"
)

func Get(name string) (types.Language, error) {
	implementation, err := getImplementation(name)
	if err != nil {
		return nil, err
	}

	return base.New(implementation), nil
}

func getImplementation(name string) (implementation.Implementation, error) {
	switch name {
	case "ruby":
		return ruby.Get(), nil
	default:
		return nil, fmt.Errorf("unsupported language '%s'", name)
	}
}
