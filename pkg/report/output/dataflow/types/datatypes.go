package types

type Datatype struct {
	UUID         string             `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	CategoryUUID string             `json:"category_uuid,omitempty" yaml:"category_uuid,omitempty"`
	Name         string             `json:"name" yaml:"name"`
	Detectors    []DatatypeDetector `json:"detectors" yaml:"detectors"`
}

type DatatypeDetector struct {
	Name      string             `json:"name" yaml:"name"`
	Locations []DatatypeLocation `json:"locations" yaml:"locations"`
}

type DatatypeLocation struct {
	Filename   string               `json:"filename" yaml:"filename"`
	LineNumber int                  `json:"line_number" yaml:"line_number"`
	Encrypted  *bool                `json:"encrypted,omitempty" yaml:"encrypted,omitempty"`
	VerifiedBy []DatatypeVerifiedBy `json:"verified_by,omitempty" yaml:"verified_by,omitempty"`
}

type DatatypeVerifiedBy struct {
	Detector   string  `json:"detector" yaml:"detector"`
	Filename   *string `json:"filename,omitempty" yaml:"filename,omitempty"`
	Linenumber *int    `json:"line_number,omitempty" yaml:"line_number,omitempty"`
}
