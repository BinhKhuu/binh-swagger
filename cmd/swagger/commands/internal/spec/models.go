package spec

// ModelSpec todo remove Packagename (will be model) outputpath (will be supplied in command) output file (will be supplied in commands)
type ModelSpec struct {
	PackageName string      `yaml:"package_name"`
	Name        string      `yaml:"name"`
	Fields      []FieldSpec `yaml:"fields"`
}

type FieldSpec struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	JSON string `yaml:"json,omitempty"`
}

type RouteSpec struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	Model  string `yaml:"model"`
}
