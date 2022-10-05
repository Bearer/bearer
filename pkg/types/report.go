package types

type Report struct {
	Path     string
	Artifact Artifact
	Metadata Metadata
}

type Metadata struct {
	Version string
}

func (report *Report) Failed() bool {
	return false
}
