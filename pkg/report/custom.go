package report

import (
	"github.com/bearer/curio/pkg/report/detectors"
	"github.com/bearer/curio/pkg/report/schema"
	"github.com/bearer/curio/pkg/report/source"
	"github.com/rs/zerolog/log"
	"github.com/wlredeye/jsonlines"
)

type CustomReport interface {
	AddCustomDetection(ruleName string, source source.Source, value schema.Schema)
}

type CustomDetection struct {
	Type         DetectionType  `json:"type"`
	DetectorType detectors.Type `json:"detector_type"`
	CommitSHA    string         `json:"commit_sha"`
	Source       source.Source  `json:"source"`
	Value        schema.Schema  `json:"value"`
}

func (report *JsonLinesReport) AddCustomDetection(ruleName string, source source.Source, value schema.Schema) {
	data := &CustomDetection{
		Type:         TypeCustom,
		DetectorType: detectors.Type(ruleName),
		Source:       source,
		Value:        value,
	}
	if data.Source.LineNumber != nil {
		data.CommitSHA = report.Blamer.SHAForLine(data.Source.Filename, *data.Source.LineNumber)
	}

	detectionsToAdd := []*CustomDetection{data}

	err := jsonlines.Encode(report.File, &detectionsToAdd)
	if err != nil {
		log.Printf("failed to encode data line %e", err)
	}
}
