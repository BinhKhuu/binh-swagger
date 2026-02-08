package spec

type ModelSpec struct {
	Name   string      `yaml:"name"`
	Fields []FieldSpec `yaml:"fields"`
}
type FieldSpec struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	JSON string `yaml:"json,omitempty"`
}

// todo check if this is still used and move this to a separate file
type RouteSpec struct {
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	Model  string `yaml:"model"`
}
