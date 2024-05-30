package pluralize

import "github.com/gertd/go-pluralize"

var pluralizer = pluralize.NewClient()

func Singular(word string) string {
	return pluralizer.Singular(word)
}
