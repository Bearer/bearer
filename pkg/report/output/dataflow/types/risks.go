package types

import "github.com/bearer/bearer/pkg/report/schema"

type RiskDetector struct {
	DetectorID string         `json:"detector_id" yaml:"detector_id"`
	Locations  []RiskLocation `json:"locations" yaml:"locations"`
}

type RiskLocation struct {
	Filename           string                 `json:"filename" yaml:"filename"`
	LineNumber         int                    `json:"line_number" yaml:"line_number"`
	DataTypeCategories []RiskDatatypeCategory `json:"data_type_categories,omitempty" yaml:"data_type_categories,omitempty"`
}

type RiskDatatypeCategory struct {
	Name       string         `json:"name,omitempty" yaml:"name,omitempty"`
	DataTypes  []RiskDatatype `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	IsPresence bool           `json:"is_presence,omitempty" yaml:"is_presence,omitempty"`
}

type RiskDatatype struct {
	Parent      *schema.Parent `json:"parent,omitempty" yaml:"parent,omitempty"`
	Content     string         `json:"content,omitempty"  yaml:"content,omitempty"`
	FieldName   string         `json:"field_name,omitempty" yaml:"field_name,omitempty"`
	ObjectName  string         `json:"object_name,omitempty" yaml:"object_name,omitempty"`
	SubjectName *string        `json:"subject_name,omitempty" yaml:"subject_name,omitempty"`
	Stored      bool           `json:"stored" yaml:"stored"`
}
