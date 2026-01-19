package spec

type ModelSpec struct {
	PackageName string      `yaml:"package_name"`
	OutputPath  string      `yaml:"output_path"`
	OutputFile  string      `yaml:"output_file"`
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
