package ruby

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"

	"github.com/bearer/bearer/internal/classification/schema"
	"github.com/bearer/bearer/internal/report/detectors"
	"github.com/bearer/bearer/internal/scanner/ast/query"
	"github.com/bearer/bearer/internal/scanner/ast/tree"
	detectortypes "github.com/bearer/bearer/internal/scanner/detectors/types"

	"github.com/bearer/bearer/internal/languages/ruby/analyzer"
	"github.com/bearer/bearer/internal/languages/ruby/detectors/object"
	stringdetector "github.com/bearer/bearer/internal/languages/ruby/detectors/string"
	"github.com/bearer/bearer/internal/languages/ruby/pattern"
	"github.com/bearer/bearer/internal/scanner/detectors/datatype"
	"github.com/bearer/bearer/internal/scanner/detectors/insecureurl"
	"github.com/bearer/bearer/internal/scanner/detectors/stringliteral"
	"github.com/bearer/bearer/internal/scanner/language"
)

type implementation struct {
	pattern pattern.Pattern
}

func Get() language.Language {
	return &implementation{}
}

func (*implementation) Name() string {
	return "ruby"
}

func (*implementation) EnryLanguages() []string {
	return []string{"Ruby"}
}

func (*implementation) NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector {
	return []detectortypes.Detector{
		object.New(querySet),
		datatype.New(detectors.DetectorRuby, schemaClassifier),
		stringdetector.New(querySet),
		stringliteral.New(querySet),
		insecureurl.New(querySet),
	}
}

func (*implementation) SitterLanguage() *sitter.Language {
	return ruby.GetLanguage()
}

func (language *implementation) Pattern() language.Pattern {
	return &language.pattern
}

func (*implementation) NewAnalyzer(builder *tree.Builder) language.Analyzer {
	return analyzer.New(builder)
}
