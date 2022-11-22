package types

type Component struct {
	Name      string              `json:"name"`
	UUID      string              `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Locations []ComponentLocation `json:"locations"`
}

type ComponentLocation struct {
	Detector   string `json:"detector"`
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
}
