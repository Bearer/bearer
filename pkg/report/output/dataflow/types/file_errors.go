package types

type Error struct {
	Type     string `json:"type" yaml:"type"`
	Filename string `json:"filename" yaml:"filename"`
	Error    string `json:"error" yaml:"filename"`
}
