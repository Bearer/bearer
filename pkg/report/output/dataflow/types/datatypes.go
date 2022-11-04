package types

type Datatype struct {
	Name      string             `json:"name"`
	Detectors []DatatypeDetector `json:"detectors"`
}

type DatatypeDetector struct {
	Name      string             `json:"name"`
	Stored    bool               `json:"stored"`
	Locations []DatatypeLocation `json:"locations"`
}

type DatatypeLocation struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
}
