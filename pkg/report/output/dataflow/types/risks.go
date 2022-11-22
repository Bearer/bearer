package types

type RiskDetector struct {
	DetectorID string         `json:"detector_id" yaml:"detector_id"`
	DataTypes  []RiskDatatype `json:"data_types" yaml:"data_types"`
}

type RiskDatatype struct {
	Name         string         `json:"name" yaml:"name"`
	UUID         string         `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	CategoryUUID string         `json:"category_uuid,omitempty" yaml:"category_uuid,omitempty"`
	Stored       bool           `json:"stored" yaml:"stored"`
	Locations    []RiskLocation `json:"locations" yaml:"locations"`
}

type RiskLocation struct {
	Filename   string `json:"filename" yaml:"filename"`
	LineNumber int    `json:"line_number" yaml:"line_number"`
}
