package types

type Component struct {
	Name      string              `json:"name" yaml:"name"`
	Type      string              `json:"type" yaml:"type"`
	SubType   string              `json:"sub_type" yaml:"sub_type"`
	UUID      string              `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Locations []ComponentLocation `json:"locations" yaml:"locations"`
}

type Dependency struct {
	Name             string `json:"name" yaml:"name"`
	Version          string `json:"version" yaml:"version"`
	Filename         string `json:"filename" yaml:"filename"`
	Detector         string `json:"detector" yaml:"detector"`
	DetectorLanguage string `json:"-" yaml:"-"`
}

type ComponentLocation struct {
	Detector     string `json:"detector" yaml:"detector"`
	FullFilename string `json:"full_filename" yaml:"full_filename"`
	Filename     string `json:"filename" yaml:"filename"`
	LineNumber   int    `json:"line_number" yaml:"line_number"`
}
