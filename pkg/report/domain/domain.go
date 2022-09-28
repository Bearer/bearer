package domain

type Domain struct {
	Domain  string `json:"domain"`
	Context string `json:"context,omitempty"`
}
