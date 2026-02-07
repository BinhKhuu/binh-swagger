package spec

type APIConfig struct {
	Version string               `yaml:"version"`
	Models  map[string]ModelSpec `yaml:"models"`
	Paths   map[string]PathSpec  `yaml:"paths"`
}
