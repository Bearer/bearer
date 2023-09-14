package types

import "github.com/bearer/bearer/internal/report/schema"

type RiskDetector struct {
	DetectorID string         `json:"detector_id" yaml:"detector_id"`
	Locations  []RiskLocation `json:"locations" yaml:"locations"`
}

type RiskLocation struct {
	FullFilename      string         `json:"full_filename" yaml:"full_filename"`
	Filename          string         `json:"filename" yaml:"filename"`
	StartLineNumber   int            `json:"start_line_number" yaml:"start_line_number"`
	StartColumnNumber int            `json:"start_column_number" yaml:"start_column_number"`
	EndLineNumber     int            `json:"end_line_number" yaml:"end_line_number"`
	EndColumnNumber   int            `json:"end_column_number" yaml:"end_column_number"`
	Source            *schema.Source `json:"source,omitempty" yaml:"source,omitempty"`
	DataTypes         []RiskDatatype `json:"data_types,omitempty" yaml:"data_types,omitempty"`
	PresenceMatches   []RiskPresence `json:"presence_matches,omitempty" yaml:"presence_matches,omitempty"`
}

type RiskDatatype struct {
	Name         string       `json:"name,omitempty" yaml:"name,omitempty"`
	CategoryUUID string       `json:"category_uuid,omitempty" yaml:"category_uuid,omitempty"`
	Schemas      []RiskSchema `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

type RiskSchema struct {
	FieldName   string  `json:"field_name,omitempty" yaml:"field_name,omitempty"`
	ObjectName  string  `json:"object_name,omitempty" yaml:"object_name,omitempty"`
	SubjectName *string `json:"subject_name,omitempty" yaml:"subject_name,omitempty"`
}

type RiskPresence struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}
