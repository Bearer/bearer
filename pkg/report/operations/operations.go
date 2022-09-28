package operations

const (
	TypeGet    = "GET"
	TypePost   = "POST"
	TypePut    = "PUT"
	TypeDelete = "DELETE"
	TypeOther  = "OTHER"
)

type Operation struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Urls []Url  `json:"url"`
}

type Url struct {
	Url       string     `json:"url"`
	Variables []Variable `json:"variables"`
}

type Variable struct {
	Name   string
	Values []string
}
