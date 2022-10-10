package depsbase

// Dependency is a dependency that keeps the name and version
type Dependency struct {
	Group   string `json:"group"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Line    int64  `json:"lineNumber,omitempty"`
	Column  int64  `json:"columnNumber,omitempty"`
}

// DiscoveredDependency holds a list of dependencies defined in package file
type DiscoveredDependency struct {
	Provider       string       `json:"provider"`
	Language       string       `json:"language"`
	PackageManager string       `json:"package_manager`
	Dependencies   []Dependency `json:"dependencies"`
}
