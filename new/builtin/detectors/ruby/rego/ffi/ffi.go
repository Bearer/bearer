package ffi

import (
	"log"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"

	"github.com/bearer/curio/new/language"
	languagetypes "github.com/bearer/curio/new/language/types"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
)

var (
	opaqueType = types.NewObject([]*types.StaticProperty{
		{Key: "type", Value: types.S},
		{Key: "id", Value: types.N},
	}, nil)

	matchType = types.NewObject(nil, types.NewDynamicProperty(types.S, opaqueType))

	detectionType = types.NewObject([]*types.StaticProperty{
		{Key: "match_node", Value: opaqueType},
		{Key: "data", Value: types.NewObject(nil, types.NewDynamicProperty(types.S, types.A))},
	}, nil)
)

type EvaluationContext struct {
	evaluator treeevaluatortypes.Evaluator
	lang      languagetypes.Language
	nodeMap   opaqueTypeMap[language.Node]
	queryMap  opaqueTypeMap[language.Query]
}

func NewEvalContext(lang languagetypes.Language) *EvaluationContext {
	return &EvaluationContext{
		lang:     lang,
		nodeMap:  newOpaqueTypeMap[language.Node]("node"),
		queryMap: newOpaqueTypeMap[language.Query]("query"),
	}
}

func (evalContext *EvaluationContext) NewFile(evaluator treeevaluatortypes.Evaluator) {
	evalContext.evaluator = evaluator
	evalContext.nodeMap = newOpaqueTypeMap[language.Node]("node")
}

func (evalContext *EvaluationContext) NodeToRego(node *language.Node) *ast.Term {
	return evalContext.nodeMap.CastToRego(node)
}

func (evalContext *EvaluationContext) Close() {
	for _, query := range evalContext.queryMap.idToInstance {
		query.Close()
	}
}

// curio.evaluator.detections_at
func EvaluatorDetectionsAt(evalContext *EvaluationContext) func(*rego.Rego) {
	return rego.Function2(
		&rego.Function{
			Name: "curio.evaluator.detections_at",
			Decl: types.NewFunction(
				[]types.Type{opaqueType, types.S},
				types.NewArray([]types.Type{detectionType}, nil),
			),
			Memoize: true,
		},
		func(bctx rego.BuiltinContext, nodeTerm, detectorTypeTerm *ast.Term) (*ast.Term, error) {
			log.Printf("curio.evaluator.detections_at: enter")

			node, err := evalContext.nodeMap.CastToGo(nodeTerm)
			if err != nil {
				log.Printf("curio.evaluator.detections_at: err: %s", err)
				return nil, err
			}

			detectorType := string(detectorTypeTerm.Value.(ast.String))

			log.Printf("curio.evaluator.detections_at: detectorType: %s", detectorType)

			detections, err := evalContext.evaluator.NodeDetections(node, detectorType)
			if err != nil {
				return nil, err
			}

			log.Printf("curio.evaluator.detections_at: exit")

			return ast.NewTerm(detections), err
		},
	)
}

// curio.node.content
func NodeContent(evalContext *EvaluationContext) func(*rego.Rego) {
	return rego.Function1(
		&rego.Function{
			Name: "curio.node.content",
			Decl: types.NewFunction(
				[]types.Type{opaqueType},
				types.S,
			),
			Memoize: true,
		},
		func(bctx rego.BuiltinContext, nodeTerm *ast.Term) (*ast.Term, error) {
			log.Printf("curio.node.content: enter")

			node, err := evalContext.nodeMap.CastToGo(nodeTerm)
			if err != nil {
				log.Printf("curio.node.content: err: %s", err)
				return nil, err
			}

			log.Printf("curio.node.content: exit")
			return ast.StringTerm(node.Content()), nil
		},
	)
}

// curio.language.compile_sitter_query
func LangaugeCompileSitterQuery(evalContext *EvaluationContext) func(*rego.Rego) {
	return rego.Function1(
		&rego.Function{
			Name: "curio.language.compile_sitter_query",
			Decl: types.NewFunction(
				[]types.Type{types.S},
				opaqueType,
			),
			Memoize: true,
		},
		func(bctx rego.BuiltinContext, inputTerm *ast.Term) (*ast.Term, error) {
			log.Printf("curio.language.compile_sitter_query: enter")

			input := string(inputTerm.Value.(ast.String))

			query, err := evalContext.lang.CompileQuery(input)
			if err != nil {
				return nil, err
			}

			log.Printf("curio.language.compile_sitter_query: exit")
			return evalContext.queryMap.CastToRego(query), nil
		},
	)
}

// curio.query.match_at
func QueryMatchAt(evalContext *EvaluationContext) func(*rego.Rego) {
	return rego.Function2(
		&rego.Function{
			Name: "curio.query.match_at",
			Decl: types.NewFunction(
				[]types.Type{opaqueType, opaqueType},
				types.NewArray([]types.Type{matchType}, nil),
			),
			Memoize: true,
		},
		func(bctx rego.BuiltinContext, queryTerm, nodeTerm *ast.Term) (*ast.Term, error) {
			log.Printf("curio.query.match_at: enter")

			query, err := evalContext.queryMap.CastToGo(queryTerm)
			if err != nil {
				log.Printf("curio.query.match_at: err: %s", err)
				return nil, err
			}

			node, err := evalContext.nodeMap.CastToGo(nodeTerm)
			if err != nil {
				log.Printf("curio.query.match_at: err: %s", err)
				return nil, err
			}

			results, err := query.MatchAt(node)
			if err != nil {
				log.Printf("curio.query.match_at: err: %s", err)
				return nil, err
			}

			resultTerms := make([]*ast.Term, len(results))

			for i, result := range results {
				resultTerms[i] = translateQueryResult(evalContext, result)
			}

			log.Printf("curio.query.match_at: exit")
			return ast.ArrayTerm(resultTerms...), nil
		},
	)
}

// curio.query.match_once_at
func QueryMatchOnceAt(evalContext *EvaluationContext) func(*rego.Rego) {
	return rego.Function2(
		&rego.Function{
			Name: "curio.query.match_once_at",
			Decl: types.NewFunction(
				[]types.Type{opaqueType, opaqueType},
				matchType,
			),
			Memoize: true,
		},
		func(bctx rego.BuiltinContext, queryTerm, nodeTerm *ast.Term) (*ast.Term, error) {
			log.Printf("curio.query.match_once_at: enter")

			query, err := evalContext.queryMap.CastToGo(queryTerm)
			if err != nil {
				log.Printf("curio.query.match_once_at: err: %s", err)
				return nil, err
			}

			node, err := evalContext.nodeMap.CastToGo(nodeTerm)
			if err != nil {
				log.Printf("curio.query.match_once_at: err: %s", err)
				return nil, err
			}

			result, err := query.MatchOnceAt(node)
			if err != nil {
				log.Printf("curio.query.match_once_at: err: %s", err)
				return nil, err
			}

			log.Printf("curio.query.match_once_at: exit")
			return translateQueryResult(evalContext, result), nil
		},
	)
}

func translateQueryResult(evalContext *EvaluationContext, result language.QueryResult) *ast.Term {
	resultTerms := make([][2]*ast.Term, len(result))
	i := 0

	for name, node := range result {
		resultTerms[i] = ast.Item(ast.StringTerm(name), evalContext.nodeMap.CastToRego(node))
		i += 1
	}

	return ast.ObjectTerm(resultTerms...)
}
