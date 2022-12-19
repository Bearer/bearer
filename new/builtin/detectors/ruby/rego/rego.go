package rego

import (
	"context"
	"embed"
	"fmt"
	"log"
	"strings"

	"github.com/bearer/curio/new/builtin/detectors/ruby/rego/ffi"
	detectiontypes "github.com/bearer/curio/new/detection/types"
	"github.com/bearer/curio/new/detector"
	"github.com/bearer/curio/new/language"
	languagetypes "github.com/bearer/curio/new/language/types"
	treeevaluatortypes "github.com/bearer/curio/new/treeevaluator/types"
	"github.com/open-policy-agent/opa/rego"
)

//go:embed detectors/*
var detectorsFs embed.FS

type Data struct {
	Name string
}

type regoDetector struct {
	detectorType string
	regoInstance *rego.Rego
	regoQuery    rego.PreparedEvalQuery
}

func New(lang languagetypes.Language, detectorType string) (detector.Detector, error) {
	filename := "detectors/" + detectorType + ".rego"

	detectorContent, err := detectorsFs.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	regoInstance := rego.New(
		rego.Trace(true),
		rego.Query("detections = data.curio.detectors."+detectorType+".detect_at"),
		rego.Module(filename, string(detectorContent)),
		ffi.EvaluatorNodeDetections(),
	)

	regoQuery, err := regoInstance.PrepareForEval(context.TODO())
	if err != nil {
		return nil, err
	}

	return &regoDetector{
		detectorType: detectorType,
		regoInstance: regoInstance,
		regoQuery:    regoQuery,
	}, nil
}

func (detector *regoDetector) Type() string {
	return detector.detectorType
}

func (detector *regoDetector) DetectAt(
	node *language.Node,
	evaluator treeevaluatortypes.Evaluator,
) ([]*detectiontypes.Detection, error) {
	ffiData, err := ffi.GetData()
	if err != nil {
		return nil, err
	}

	nodeInput := ffiData.NodeToRegoInput(node)

	log.Printf("REGO! detect")

	resultSet, err := detector.regoQuery.Eval(context.TODO(), rego.EvalInput(nodeInput))
	if err != nil {
		return nil, err
	}

	if len(resultSet) == 0 {
		return nil, fmt.Errorf("missing rule detect_at in '%s' detector", detector.detectorType)
	}

	if len(resultSet) != 1 {
		return nil, fmt.Errorf("expected single result from query got %d results %#v", len(resultSet), resultSet)
	}

	detections := resultSet[0].Bindings["detections"]
	log.Printf("detections: %#v", detections)

	builder := strings.Builder{}
	rego.PrintTraceWithLocation(&builder, detector.regoInstance)
	log.Printf("t: %#v", builder.String())

	return nil, nil
}

func (detector *regoDetector) Close() {}
