package interfaces

type ClassifiedFramework struct {
	data           interface{}
	Classification Classification `json:"classification"`
}

type Classification struct {
}
