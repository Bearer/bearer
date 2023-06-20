package types

type Component struct {
	Name      string              `json:"name" yaml:"name"`
	Type      string              `json:"type" yaml:"type"`
	SubType   string              `json:"sub_type" yaml:"sub_type"`
	UUID      string              `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Locations []ComponentLocation `json:"locations" yaml:"locations"`
}

type Dependency struct {
	Name     string `json:"name" yaml:"name"`
	Version  string `json:"version" yaml:"version"`
	Filename string `json:"filename" yaml:"filename"`
	// Version  Version `json:"version" yaml:"version"`
	// Detector string `json:"detector" yaml:"detector"`
}

type ComponentLocation struct {
	Detector     string `json:"detector" yaml:"detector"`
	FullFilename string `json:"full_filename" yaml:"full_filename"`
	Filename     string `json:"filename" yaml:"filename"`
	LineNumber   int    `json:"line_number" yaml:"line_number"`
}

type Version struct {
	Major string `json:"major" yaml:"major"`
	Minor string `json:"minor" yaml:"minor"`
	Patch string `json:"patch" yaml:"patch"`
}
