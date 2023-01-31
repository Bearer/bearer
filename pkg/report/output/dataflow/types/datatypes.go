package types

import (
	"github.com/bearer/curio/pkg/report/schema"
)

type Datatype struct {
	UUID         string             `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	CategoryUUID string             `json:"category_uuid,omitempty" yaml:"category_uuid,omitempty"`
	Name         string             `json:"name" yaml:"name"`
	Detectors    []DatatypeDetector `json:"detectors" yaml:"detectors"`
}

type DatatypeDetector struct {
	Name      string             `json:"name" yaml:"name"`
	Locations []DatatypeLocation `json:"locations" yaml:"locations"`
	Parent    *schema.Parent     `json:"parent,omitempty" yaml:"parent,omitempty"`
}

type DatatypeLocation struct {
	Filename    string               `json:"filename" yaml:"filename"`
	LineNumber  int                  `json:"line_number" yaml:"line_number"`
	Encrypted   *bool                `json:"encrypted,omitempty" yaml:"encrypted,omitempty"`
	VerifiedBy  []DatatypeVerifiedBy `json:"verified_by,omitempty" yaml:"verified_by,omitempty"`
	Stored      *bool                `json:"stored,omitempty" yaml:"stored,omitempty"`
	Parent      *schema.Parent       `json:"parent,omitempty" yaml:"parent,omitempty"`
	FieldName   string               `json:"field_name,omitempty" yaml:"field_name,omitempty"`
	ObjectName  string               `json:"object_name,omitempty" yaml:"object_name,omitempty"`
	SubjectName *string              `json:"subject_name,omitempty" yaml:"subject_name,omitempty"`
}

type DatatypeVerifiedBy struct {
	Detector   string  `json:"detector" yaml:"detector"`
	Filename   *string `json:"filename,omitempty" yaml:"filename,omitempty"`
	Linenumber *int    `json:"line_number,omitempty" yaml:"line_number,omitempty"`
}
