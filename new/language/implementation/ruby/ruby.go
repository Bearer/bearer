package ruby

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"

	detectortypes "github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/report/detectors"

	"github.com/bearer/bearer/new/detector/implementation/generic/datatype"
	"github.com/bearer/bearer/new/detector/implementation/generic/insecureurl"
	"github.com/bearer/bearer/new/detector/implementation/generic/stringliteral"
	"github.com/bearer/bearer/new/detector/implementation/ruby/object"
	stringdetector "github.com/bearer/bearer/new/detector/implementation/ruby/string"
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/implementation/ruby/analyzer"
	"github.com/bearer/bearer/new/language/implementation/ruby/pattern"
)

type rubyImplementation struct {
	pattern pattern.Pattern
}

func Get() implementation.Implementation {
	return &rubyImplementation{}
}

func (*rubyImplementation) Name() string {
	return "ruby"
}

func (*rubyImplementation) EnryLanguages() []string {
	return []string{"Ruby"}
}

func (*rubyImplementation) NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector {
	return []detectortypes.Detector{
		object.New(querySet),
		datatype.New(detectors.DetectorRuby, schemaClassifier),
		stringdetector.New(querySet),
		stringliteral.New(querySet),
		insecureurl.New(querySet),
	}
}

func (*rubyImplementation) SitterLanguage() *sitter.Language {
	return ruby.GetLanguage()
}

func (implementation *rubyImplementation) Pattern() implementation.Pattern {
	return &implementation.pattern
}

func (*rubyImplementation) NewAnalyzer(builder *tree.Builder) implementation.Analyzer {
	return analyzer.New(builder)
}
