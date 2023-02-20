package types

import "github.com/bearer/bearer/pkg/report/schema"

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
}

type RiskLocation struct {
	Filename    string         `json:"filename" yaml:"filename"`
	LineNumber  int            `json:"line_number" yaml:"line_number"`
	Parent      *schema.Parent `json:"parent,omitempty" yaml:"parent,omitempty"`
	FieldName   string         `json:"field_name,omitempty" yaml:"field_name,omitempty"`
	ObjectName  string         `json:"object_name,omitempty" yaml:"object_name,omitempty"`
	SubjectName *string        `json:"subject_name,omitempty" yaml:"subject_name,omitempty"`
}

type RiskDetectionLocation struct {
	*RiskLocation `json:",inline" yaml:",inline"`
	Content       string `json:"content" yaml:"content"`
}

type RiskDetection struct {
	DetectorID string                  `json:"detector_id" yaml:"detector_id"`
	Locations  []RiskDetectionLocation `json:"locations" yaml:"locations"`
}
