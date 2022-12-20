package rego

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bearer/curio/new/builtin/detectors/ruby/rego/ffi"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/topdown"
)

//go:embed detectors/*
var detectorsFs embed.FS

type Data struct {
	Name string
}

type regoDetector struct {
	evaluationContext *ffi.EvaluationContext
	detectorType      string
	regoInstance      *rego.Rego
	regoQuery         rego.PreparedEvalQuery
}

func New(evaluationContext *ffi.EvaluationContext, detectorType string) (detector.Detector, error) {
	filename := "detectors/" + detectorType + ".rego"

	detectorContent, err := detectorsFs.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	regoInstance := rego.New(
		rego.Query("detections = data.curio.detectors."+detectorType+".detections_at_input"),
		rego.Module(filename, string(detectorContent)),
		ffi.EvaluatorDetectionsAt(evaluationContext),
		ffi.NodeContent(evaluationContext),
		ffi.LangaugeCompileSitterQuery(evaluationContext),
		ffi.QueryMatchAt(evaluationContext),
		ffi.QueryMatchOnceAt(evaluationContext),
	)

	regoQuery, err := regoInstance.PrepareForEval(context.TODO())
	if err != nil {
		return nil, err
	}

	return &regoDetector{
		evaluationContext: evaluationContext,
		detectorType:      detectorType,
		regoInstance:      regoInstance,
		regoQuery:         regoQuery,
	}, nil
}

func (detector *regoDetector) Type() string {
	return detector.detectorType
}

func (detector *regoDetector) DetectAt(
	node *language.Node,
	evaluator treeevaluatortypes.Evaluator,
) (*ast.Array, error) {
	log.Printf("REGO! detect")

	tracer := topdown.NewBufferTracer()
	resultSet, err := detector.regoQuery.Eval(
		context.TODO(),
		rego.EvalParsedInput(detector.evaluationContext.NodeToRego(node).Value),
	)
	if err != nil {
		return nil, err
	}

	if len(resultSet) == 0 {
		return nil, fmt.Errorf("missing rule detections_at in '%s' detector", detector.detectorType)
	}

	if len(resultSet) != 1 {
		return nil, fmt.Errorf("expected single result from query got %d results %#v", len(resultSet), resultSet)
	}

	detections := resultSet[0].Bindings["detections"]
	log.Printf("detections: %#v", detections)

	builder := strings.Builder{}
	topdown.PrettyTraceWithLocation(&builder, *tracer)
	log.Printf("t: %s", builder.String())

	detectionsValue, err := ast.InterfaceToValue(detections)
	if err != nil {
		return nil, err
	}

	detectionsArray, ok := detectionsValue.(*ast.Array)
	if !ok {
		return nil, errors.New("expected array but got other term")
	}

	return detectionsArray, nil
}

func (detector *regoDetector) Close() {}
