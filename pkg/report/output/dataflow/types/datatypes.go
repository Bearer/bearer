package types

type Datatype struct {
	Name      string             `json:"name" yaml:"name"`
	UUID      string             `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Detectors []DatatypeDetector `json:"detectors" yaml:"detectors"`
}

type DatatypeDetector struct {
	Name      string             `json:"name" yaml:"name"`
	Locations []DatatypeLocation `json:"locations" yaml:"locations"`
}

type DatatypeLocation struct {
	Filename   string `json:"filename" yaml:"filename"`
	LineNumber int    `json:"line_number" yaml:"line_number"`
}
