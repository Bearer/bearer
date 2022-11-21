package types

type RiskDetector struct {
	DetectorID string         `json:"detector_id"`
	DataTypes  []RiskDatatype `json:"data_types"`
}

type RiskDatatype struct {
	Name         string         `json:"name"`
	UUID         string         `json:"uuid,omitempty"`
	CategoryUUID string         `json:"category_uuid,omitempty"`
	Stored       bool           `json:"stored"`
	Locations    []RiskLocation `json:"locations"`
}

type RiskLocation struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
}
