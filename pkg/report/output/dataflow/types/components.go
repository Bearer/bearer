package types

type Component struct {
	Name      string
	Locations []ComponentLocation
}

type ComponentLocation struct {
	Detector   string
	Filename   string
	LineNumber int
}
