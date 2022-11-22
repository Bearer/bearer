package types

type Component struct {
	Name      string
	UUID      string              `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Locations []ComponentLocation
}

type ComponentLocation struct {
	Detector   string
	Filename   string
	LineNumber int
}
