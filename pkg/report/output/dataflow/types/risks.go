package types

// - detector_id: "rails_logger_detector"
// data_types:
//   - name: "email"
// 	locations:
// 	  - filename: "app/models/user.rb"
// 		line_number: 5
// 	  - filename: "app/models/employee.rb"
// 		line_number: 5

type RiskDetector struct {
	DetectorID string         `json:"detector_id"`
	DataTypes  []RiskDatatype `json:"data_types"`
}

type RiskDatatype struct {
	Name      string         `json:"name"`
	Locations []RiskLocation `json:"locations"`
}

type RiskLocation struct {
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
}
