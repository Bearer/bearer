package pluralize

import "github.com/gertd/go-pluralize"

var pluralizer = pluralize.NewClient()

func Singular(word string) string {
	return pluralizer.Singular(word)
}

func Plural(word string) string {
	return pluralizer.Plural(word)
}
