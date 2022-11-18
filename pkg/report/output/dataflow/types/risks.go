package types

type RiskDetector struct {
	DetectorID string         `json:"detector_id"`
	DataTypes  []RiskDatatype `json:"data_types"`
}

type RiskDatatype struct {
	Name         string         `json:"name"`
	UUID         string         `json:"uuid"`
	CategoryUUID string         `json:"category_uuid"`
	Stored       bool           `json:"stored"`
	Locations    []RiskLocation `json:"locations"`
}

type RiskLocation struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
}
