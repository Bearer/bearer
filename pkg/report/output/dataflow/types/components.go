package types

type Component struct {
	Name      string              `json:"name" yaml:"name"`
	Type      string              `json:"type" yaml:"type"`
	UUID      string              `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Locations []ComponentLocation `json:"locations" yaml:"locations"`
}

type ComponentLocation struct {
	Detector   string `json:"detector" yaml:"detector"`
	Filename   string `json:"filename" yaml:"filename"`
	LineNumber int    `json:"line_number" yaml:"line_number"`
}
