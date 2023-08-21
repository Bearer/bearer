package java

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"

	detectortypes "github.com/bearer/bearer/new/detector/types"
	"github.com/bearer/bearer/pkg/ast/query"
	"github.com/bearer/bearer/pkg/ast/tree"
	"github.com/bearer/bearer/pkg/classification/schema"
	"github.com/bearer/bearer/pkg/report/detectors"

	"github.com/bearer/bearer/new/detector/implementation/generic/datatype"
	"github.com/bearer/bearer/new/detector/implementation/generic/insecureurl"
	"github.com/bearer/bearer/new/detector/implementation/generic/stringliteral"
	"github.com/bearer/bearer/new/detector/implementation/java/object"
	stringdetector "github.com/bearer/bearer/new/detector/implementation/java/string"
	"github.com/bearer/bearer/new/language/implementation"
	"github.com/bearer/bearer/new/language/implementation/java/analyzer"
	"github.com/bearer/bearer/new/language/implementation/java/pattern"
)

type javaImplementation struct {
	pattern pattern.Pattern
}

func Get() implementation.Implementation {
	return &javaImplementation{}
}

func (*javaImplementation) Name() string {
	return "java"
}

func (*javaImplementation) EnryLanguages() []string {
	return []string{"Java"}
}

func (*javaImplementation) NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector {
	return []detectortypes.Detector{
		object.New(querySet),
		datatype.New(detectors.DetectorJava, schemaClassifier),
		stringdetector.New(querySet),
		stringliteral.New(querySet),
		insecureurl.New(querySet),
	}
}

func (*javaImplementation) SitterLanguage() *sitter.Language {
	return java.GetLanguage()
}

func (implementation *javaImplementation) Pattern() implementation.Pattern {
	return &implementation.pattern
}

func (*javaImplementation) NewAnalyzer(builder *tree.Builder) implementation.Analyzer {
	return analyzer.New(builder)
}
