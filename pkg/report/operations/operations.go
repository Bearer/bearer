package operations

const (
	TypeGet     = "GET"
	TypePost    = "POST"
	TypePut     = "PUT"
	TypeDelete  = "DELETE"
	TypePatch   = "PATCH"
	TypeHead    = "HEAD"
	TypeOptions = "OPTIONS"
	TypeConnect = "CONNECT"
	TypeTrace   = "TRACE"
)

type Operation struct {
	Path string `json:"path" yaml:"path"`
	Type string `json:"type" yaml:"type"`
	Urls []Url  `json:"url" yaml:"url"`
}

type Url struct {
	Url       string     `json:"url" yaml:"url"`
	Variables []Variable `json:"variables" yaml:"variables"`
}

type Variable struct {
	Name   string
	Values []string
}
