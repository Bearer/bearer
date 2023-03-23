package types

import "github.com/bearer/bearer/pkg/report/schema"

type RiskDetector struct {
	DetectorID string         `json:"detector_id" yaml:"detector_id"`
	Locations  []RiskLocation `json:"locations" yaml:"locations"`
}

type RiskLocation struct {
	Filename           string                 `json:"filename" yaml:"filename"`
	LineNumber         int                    `json:"line_number" yaml:"line_number"`
	DataTypeCategories []RiskDatatypeCategory `json:"categories,omitempty" yaml:"categories,omitempty"`
}

type RiskDatatypeCategory struct {
	Category  string         `json:"category,omitempty" yaml:"category,omitempty"`
	Name      string         `json:"name,omitempty" yaml:"name,omitempty"`
	DataTypes []RiskDatatype `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	Parent    *schema.Parent `json:"parent,omitempty" yaml:"parent,omitempty"`
	Stored    *bool          `json:"stored,omitempty" yaml:"stored,omitempty"`
}

type RiskDatatype struct {
	Parent      *schema.Parent `json:"parent,omitempty" yaml:"parent,omitempty"`
	FieldName   string         `json:"field_name,omitempty" yaml:"field_name,omitempty"`
	ObjectName  string         `json:"object_name,omitempty" yaml:"object_name,omitempty"`
	SubjectName *string        `json:"subject_name,omitempty" yaml:"subject_name,omitempty"`
	Stored      bool           `json:"stored" yaml:"stored"`
}
