package reviewdog

// Based on schema here:
// https://raw.githubusercontent.com/reviewdog/reviewdog/master/proto/rdf/jsonschema/DiagnosticResult.jsonschema
// Not all keys are implmented as not all are relevant to Bearer

// not yet implmented
type Suggestion struct{}

type LocationPosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type LocationRange struct {
	Start LocationPosition `json:"start"`
	End   LocationPosition `json:"end"`
}

type Location struct {
	Path  string        `json:"path"`
	Range LocationRange `json:"range"`
}

type Code struct {
	RuleId           string `json:"value"`
	DocumentationUrl string `json:"url"`
}

type Diagnostic struct {
	Message     string       `json:"message"`
	Location    Location     `json:"location"`
	Severity    string       `json:"severity"`
	Suggestions []Suggestion `json:"suggestions"`
	Code        Code         `json:"code"`
}

type Source struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ReviewdogOutput struct {
	Source      Source       `json:"source"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}
