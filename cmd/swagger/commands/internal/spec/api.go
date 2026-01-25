package spec

type APIConfig struct {
	ProjectRoot string               `yaml:"project_root"`
	Version     string               `yaml:"version"`
	Models      map[string]ModelSpec `yaml:"models"`
	Paths       map[string]PathSpec  `yaml:"paths"`
}
