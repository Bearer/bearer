package domain

type Domain struct {
	Domain  string `json:"domain" yaml:"domain"`
	Context string `json:"context,omitempty" yaml:"context,omitempty"`
}
