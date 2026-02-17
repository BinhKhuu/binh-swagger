package spec

type APIConfig struct {
	Version string               `yaml:"version"`
	Models  map[string]ModelSpec `yaml:"definitions"`
	Paths   map[string]PathSpec  `yaml:"paths"`
}
