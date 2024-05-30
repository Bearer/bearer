package depsbase

// Dependency is a dependency that keeps the name and version
type Dependency struct {
	Group   string `json:"group" yaml:"group"`
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	Line    int64  `json:"lineNumber,omitempty" yaml:"lineNumber,omitempty"`
	Column  int64  `json:"columnNumber,omitempty" yaml:"columnNumber,omitempty"`
}

// DiscoveredDependency holds a list of dependencies defined in package file
type DiscoveredDependency struct {
	Provider       string       `json:"provider" yaml:"provider"`
	Language       string       `json:"language" yaml:"language"`
	PackageManager string       `json:"package_manager" yaml:"package_manager"`
	Dependencies   []Dependency `json:"dependencies" yaml:"dependencies"`
}
