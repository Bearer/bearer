package types

import "github.com/bearer/curio/pkg/report/schema"

type RiskDetector struct {
	DetectorID string         `json:"detector_id" yaml:"detector_id"`
	DataTypes  []RiskDatatype `json:"data_types" yaml:"data_types"`
}

type RiskDatatype struct {
	Name         string         `json:"name" yaml:"name"`
	UUID         string         `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	CategoryUUID string         `json:"category_uuid,omitempty" yaml:"category_uuid,omitempty"`
	Stored       bool           `json:"stored" yaml:"stored"`
	Locations    []RiskLocation `json:"locations" yaml:"locations"`
	Parent       *schema.Parent `json:"parent,omitempty" yaml:"parent,omitempty"`
}

type RiskLocation struct {
	Filename   string `json:"filename" yaml:"filename"`
	LineNumber int    `json:"line_number" yaml:"line_number"`
}

type RiskDetectionLocation struct {
	*RiskLocation
	Content string `json:"content" yaml:"content"`
}

type RiskDetection struct {
	DetectorID string         `json:"detector_id" yaml:"detector_id"`
	Locations  []RiskDetectionLocation `json:"locations" yaml:"locations"`
}
