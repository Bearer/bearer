package language

import (
	"fmt"

	"github.com/bearer/bearer/new/language/base"
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/implementation/java"
	"github.com/bearer/bearer/new/language/implementation/javascript"
	"github.com/bearer/bearer/new/language/implementation/ruby"
	"github.com/bearer/bearer/new/language/types"
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
	case "java":
		return java.Get(), nil
	case "ruby":
		return ruby.Get(), nil
	case "javascript":
		return javascript.Get(), nil
	default:
		return nil, fmt.Errorf("unsupported language '%s'", name)
	}
}
