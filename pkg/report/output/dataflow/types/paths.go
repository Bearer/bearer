package types

type Path struct {
	DetectorName string       `json:"detector_name" yaml:"detector_name"`
	Detections   []*Detection `json:"detections" yaml:"detections"`
}

type Detection struct {
	FullFilename string   `json:"full_filename" yaml:"full_filename"`
	FullName     string   `json:"full_name" yaml:"full_name"`
	LineNumber   *int     `json:"line_number" yaml:"line_number"`
	Path         string   `json:"path" yaml:"path"`
	HttpMethod   string   `json:"http_method" yaml:"http_method"`
	Urls         []string `json:"urls" yaml:"urls"`
}
