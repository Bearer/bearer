package languages

import (
	"github.com/bearer/bearer/pkg/languages/golang"
	"github.com/bearer/bearer/pkg/languages/java"
	"github.com/bearer/bearer/pkg/languages/javascript"
	"github.com/bearer/bearer/pkg/languages/php"
	"github.com/bearer/bearer/pkg/languages/python"
	"github.com/bearer/bearer/pkg/languages/ruby"
	"github.com/bearer/bearer/pkg/languages/rust"
	"github.com/bearer/bearer/pkg/scanner/language"
)

func Default() []language.Language {
	return []language.Language{
		golang.Get(),
		java.Get(),
		javascript.Get(),
		php.Get(),
		python.Get(),
		ruby.Get(),
		rust.Get(),
	}
}
