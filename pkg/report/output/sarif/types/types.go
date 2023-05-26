package sarif

type Description struct {
	Text string `json:"text"`
}

type Configuration struct {
	Level string `json:"level"`
}

type Properties struct {
	Tags      []string `json:"tags"` // Could add Data Types as Tags!
	Precision string   `json:"precision"`
}

type Help struct {
	Text     string `json:"text"`
	Markdown string `json:"markdown"`
}

type Rule struct {
	Id                   string        `json:"id"`
	Name                 string        `json:"name"`
	Kind                 string        `json:"kind,omitempty"`
	ShortDescription     Description   `json:"shortDescription"`
	FullDescription      Description   `json:"fullDescription"`
	DefaultConfiguration Configuration `json:"defaultConfiguration"`
	Properties           *Properties   `json:"properties,omitempty"`
	Help                 Help          `json:"help"`
}

type Driver struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

type Tool struct {
	Driver Driver `json:"driver"`
}

type Message struct {
	Text string `json:"text"`
}

type ArtifactLocation struct {
	URI string `json:"uri"`
}

type PhysicalLocation struct {
	ArtifactLocation ArtifactLocation `json:"artifactLocation"`
	Region           Region           `json:"region"`
}

type Region struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
	EndColumn   int `json:"endColumn"`
	EndLine     int `json:"endLine"`
}

type Location struct {
	PhysicalLocation PhysicalLocation `json:"physicalLocation"`
}

type PartialFingerprints struct {
	PrimaryLocationLineHash               string `json:"primaryLocationLineHash,omitempty"`
	PrimaryLocationStartColumnFingerprint string `json:"primaryLocationStartColumnFingerprint,omitempty"`
}

type Result struct {
	RuleId              string               `json:"ruleId"`
	RuleIndex           int                  `json:"ruleIndex,omitempty"`
	Message             Message              `json:"message"`
	Locations           []Location           `json:"locations"`
	PartialFingerprints *PartialFingerprints `json:"partialFingerprints,omitempty"`
}

type Run struct {
	Tool    Tool     `json:"tool"`
	Results []Result `json:"results"`
}

type SarifOutput struct {
	Schema  string `json:"$schema"`
	Version string `json:"version"`
	Runs    []Run  `json:"runs"`
}
