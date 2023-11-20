package types

type IgnoredFingerprint struct {
	Author        *string `json:"author,omitempty"`
	Comment       *string `json:"comment,omitempty"`
	FalsePositive bool    `json:"false_positive"`
	IgnoredAt     string  `json:"ignored_at"`
}

type SortedIgnoredFingerpint struct {
	FingerprintId      string
	IgnoredFingerprint IgnoredFingerprint
}
