package interfacedetector

import (
	"github.com/bearer/bearer/internal/parser"
	parsercontext "github.com/bearer/bearer/internal/parser/context"
	"github.com/bearer/bearer/internal/parser/interfaces"
	"github.com/bearer/bearer/internal/report"
	reportinterface "github.com/bearer/bearer/internal/report/interfaces"

	"github.com/bearer/bearer/internal/report/detectors"
	"github.com/bearer/bearer/internal/report/values"
)

type Request struct {
	Tree             *parser.Tree
	Report           report.Report
	DetectorType     detectors.Type
	AcceptExpression func(node *parser.Node) bool
	PathAllowed      bool
	ContextResolver  *parsercontext.Resolver
}

func Detect(req *Request) error {
	return req.Tree.WalkRootValues(func(node *parser.Node) {
		if req.AcceptExpression != nil && !req.AcceptExpression(node) {
			return
		}

		if interfaceType, isInterface := interfaces.GetType(node.Value(), req.PathAllowed); isInterface {

			value := node.Value()

			if req.ContextResolver != nil {
				value = ReplaceSimpleVariables(req.ContextResolver, node)
			}

			req.Report.AddInterface(req.DetectorType, reportinterface.Interface{
				Type:  interfaceType,
				Value: value,
			}, node.Source(true))
		}
	})
}

func ReplaceSimpleVariables(contextResolver *parsercontext.Resolver, node *parser.Node) *values.Value {
	existingValues := node.Value()
	newParts := []values.Part{}
	for _, part := range existingValues.Parts {

		if varReference, ok := part.(*values.VariableReference); ok {
			resolvedValue, err := contextResolver.VariableStringValue(node.ID(), varReference.Identifier.Name)
			if err == nil {
				newParts = append(newParts, &values.String{Type: values.PartTypeString, Value: resolvedValue})
				continue
			}
		}

		newParts = append(newParts, part)
	}

	existingValues.Parts = newParts

	return existingValues
}
