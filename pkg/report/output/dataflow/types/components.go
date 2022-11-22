package types

type Component struct {
	Name      string              `json:"name"`
	Locations []ComponentLocation `json:"locations"`
}

type ComponentLocation struct {
	Detector   string `json:"detector"`
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
}
