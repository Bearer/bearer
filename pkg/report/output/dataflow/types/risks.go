package types

import "github.com/bearer/bearer/pkg/report/schema"

type RiskDetector struct {
	DetectorID string         `json:"detector_id" yaml:"detector_id"`
	Locations  []RiskLocation `json:"locations" yaml:"locations"`
}

type RiskLocation struct {
	Filename   string         `json:"filename" yaml:"filename"`
	LineNumber int            `json:"line_number" yaml:"line_number"`
	Parent     *schema.Parent `json:"parent,omitempty" yaml:"parent,omitempty"`
	Datatypes  []RiskDatatype `json:"data_types,omitempty" yaml:"data_types,omitempty"`
}

type RiskDatatype struct {
	Content      string  `json:"content,omitempty"  yaml:"content,omitempty"`
	FieldName    string  `json:"field_name,omitempty" yaml:"field_name,omitempty"`
	ObjectName   string  `json:"object_name,omitempty" yaml:"object_name,omitempty"`
	SubjectName  *string `json:"subject_name,omitempty" yaml:"subject_name,omitempty"`
	Name         string  `json:"name,omitempty" yaml:"name,omitempty"`
	UUID         string  `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	CategoryUUID string  `json:"category_uuid,omitempty" yaml:"category_uuid,omitempty"`
	Stored       bool    `json:"stored" yaml:"stored"`
}
