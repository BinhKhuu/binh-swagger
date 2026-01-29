package spec

// PathSpec defines the operations available on a single path.
// PathSpec fields are pointers to allow omitting empty operations in YAML serialization.
// Path spec omitempty with pointers allow nil checks to see if an operation is defined.
type PathSpec struct {
	Name   string
	Get    *Operation `yaml:"get,omitempty"`
	Post   *Operation `yaml:"post,omitempty"`
	Put    *Operation `yaml:"put,omitempty"`
	Patch  *Operation `yaml:"patch,omitempty"`
	Delete *Operation `yaml:"delete,omitempty"`
}

type Operation struct {
	Summary     string               `yaml:"summary"`
	OperationID string               `yaml:"operationId"`
	Produces    []string             `yaml:"produces"`
	Responses   map[int]ResponseSpec `yaml:"responses"`
}

type ResponseSpec struct {
	Description string      `yaml:"description"`
	Schema      *SchemaSpec `yaml:"schema,omitempty"`
}

type SchemaSpec struct {
	Type string `yaml:"type,omitempty"`
	Ref  string `yaml:"$ref,omitempty"`
}
