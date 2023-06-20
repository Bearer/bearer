package types

type Vulnerability struct {
	Id                   string               `json:"id"`                 // fingerprint?
	Category             string               `json:"category,omitempty"` // sast?
	Name                 string               `json:"name"`               // The name of the vulnerability. This must not include the finding's specific information.
	Description          string               `json:"description"`
	Solution             string               `json:"solution"`
	CVE                  string               `json:"cve,omitempty"`
	Severity             string               `json:"severity"`   // Info, Unknown, Low, Medium, High, or Critical.
	Confidence           string               `json:"confidence"` // Unknown
	RawSourceCodeExtract string               `json:"raw_source_code_extract"`
	Scanner              VulnerabilityScanner `json:"scanner"`
	Location             Location             `json:"location"`
	Identifiers          []Identifier         `json:"identifiers"`
}

type Identifier struct {
	Type  string `json:"type"`  // type like cwe cve we use bearer for our rules
	Name  string `json:"name"`  // Human-readable name of the identifier.
	Value string `json:"value"` // id for so for cwe this would be 123
	Url   string `json:"url"`   // link to documenation
}

type VulnerabilityScanner struct {
	Id   string `json:"id" yaml:"id"`   // bearer
	Name string `json:"name" yaml:"id"` // Bearer
}

type Scanner struct {
	Id      string `json:"id" yaml:"id"`   // bearer
	Name    string `json:"name" yaml:"id"` // Bearer
	URL     string `json:"url" yaml:"id"`  // "https://github.com/bearer/bearer"
	Vendor  Vendor `json:"vendor"`
	Version string `json:"version"` // bearer version
}

type Vendor struct {
	Name string `json:"name"` // Bearer
}

// type Commit struct {
// 	SHA string `json:"sha"` // 0000000
// }

type Location struct {
	File      string `json:"file"`
	Startline int    `json:"start_line"`
	Endline   int    `json:"end_line"`
	// Commit Commit `json:"commit"`
}

type Analyzer struct {
	Id      string `json:"id"`   // bearer-sast
	Name    string `json:"name"` // Bearer SAST
	URL     string `json:"url"`  // https://github.com/bearer/bearer
	Vendor  Vendor `json:"vendor"`
	Version string `json:"version"` // Bearer version
}

type Scan struct {
	Analyzer  Analyzer `json:"analyzer" yaml:"analyzer"`
	Scanner   Scanner  `json:"scanner" yaml:"scanner"`
	Type      string   `json:"type" yaml:"type"` // sast
	StartTime string   `json:"start_time" yaml:"start_time"`
	EndTime   string   `json:"end_time" yaml:"end_time"`
	Status    string   `json:"status" yaml:"status"` // failure or success
}

type GitLabOutput struct {
	Schema          string          `json:"$schema" yaml:"-"` // "$schema": "https://gitlab.com/gitlab-org/security-products/security-report-schemas/-/raw/master/dist/sast-report-format.json",
	Version         string          `json:"version" yaml:"version"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities" yaml:"vulnerabilities"`
	Scan            Scan            `json:"scan" yaml:"scan"`
}
