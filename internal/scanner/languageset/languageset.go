package languageset

import (
	"github.com/bearer/bearer/internal/scanner/language"
)

var languages = []language.Language

func Register(extraLanguages []language.Language) {
	languages = append(languages, extraLanguages...)
}

func Get() []language.Language {
	return languages
}

func IsSupported(id string) bool {
	for _, language := range languages {
		if language.ID() == id {
			return true
		}
	}

	return false
}
