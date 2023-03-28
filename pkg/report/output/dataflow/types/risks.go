package types

import "github.com/bearer/bearer/pkg/report/schema"

type RiskDetector struct {
	DetectorID string         `json:"detector_id" yaml:"detector_id"`
	Locations  []RiskLocation `json:"locations" yaml:"locations"`
}

type RiskLocation struct {
	Filename        string         `json:"filename" yaml:"filename"`
	LineNumber      int            `json:"line_number" yaml:"line_number"`
	Parent          *schema.Parent `json:"parent,omitempty" yaml:"parent,omitempty"`
	DataTypes       []RiskDatatype `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	PresenceMatches []RiskPresence `json:"presence_matches,omitempty" yaml:"presence_matches,omitempty"`
}

type RiskDatatype struct {
	Name    string       `json:"name,omitempty" yaml:"name,omitempty"`
	Schemas []RiskSchema `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

type RiskSchema struct {
	FieldName   string  `json:"field_name,omitempty" yaml:"field_name,omitempty"`
	ObjectName  string  `json:"object_name,omitempty" yaml:"object_name,omitempty"`
	SubjectName *string `json:"subject_name,omitempty" yaml:"subject_name,omitempty"`
}

type RiskPresence struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}
