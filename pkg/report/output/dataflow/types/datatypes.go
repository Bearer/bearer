package types

import (
	"github.com/bearer/bearer/pkg/report/schema"
)

type Datatype struct {
	UUID           string             `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	CategoryUUID   string             `json:"category_uuid,omitempty" yaml:"category_uuid,omitempty"`
	CategoryName   string             `json:"category_name,omitempty" yaml:"category_name,omitempty"`
	CategoryGroups []string           `json:"category_groups,omitempty" yaml:"category_groups,omitempty"`
	Name           string             `json:"name" yaml:"name"`
	Detectors      []DatatypeDetector `json:"detectors" yaml:"detectors"`
}

type DatatypeDetector struct {
	Name      string             `json:"name" yaml:"name"`
	Locations []DatatypeLocation `json:"locations" yaml:"locations"`
	Source    *schema.Source     `json:"source,omitempty" yaml:"source,omitempty"`
}

type DatatypeLocation struct {
	Filename          string               `json:"filename" yaml:"filename"`
	FullFilename      string               `json:"full_filename" yaml:"full_filename"`
	StartLineNumber   int                  `json:"start_line_number" yaml:"start_line_number"`
	StartColumnNumber int                  `json:"start_column_number" yaml:"start_column_number"`
	EndColumnNumber   int                  `json:"end_column_number" yaml:"end_column_number"`
	Encrypted         *bool                `json:"encrypted,omitempty" yaml:"encrypted,omitempty"`
	VerifiedBy        []DatatypeVerifiedBy `json:"verified_by,omitempty" yaml:"verified_by,omitempty"`
	Stored            *bool                `json:"stored,omitempty" yaml:"stored,omitempty"`
	Source            *schema.Source       `json:"source,omitempty" yaml:"source,omitempty"`
	FieldName         string               `json:"field_name,omitempty" yaml:"field_name,omitempty"`
	ObjectName        string               `json:"object_name,omitempty" yaml:"object_name,omitempty"`
	SubjectName       *string              `json:"subject_name,omitempty" yaml:"subject_name,omitempty"`
}

type DatatypeVerifiedBy struct {
	Detector   string  `json:"detector" yaml:"detector"`
	Filename   *string `json:"filename,omitempty" yaml:"filename,omitempty"`
	Linenumber *int    `json:"line_number,omitempty" yaml:"line_number,omitempty"`
}
