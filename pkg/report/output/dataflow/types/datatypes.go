package types

type Datatype struct {
	Name      string             `json:"name"`
	UUID      string             `json:"uuid,omitempty"`
	Detectors []DatatypeDetector `json:"detectors"`
}

type DatatypeDetector struct {
	Name      string             `json:"name"`
	Locations []DatatypeLocation `json:"locations"`
}

type DatatypeLocation struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
}
