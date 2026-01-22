package spec

type APIConfig struct {
	ProjectRoot string      `yaml:"project_root"`
	Version     string      `yaml:"version"`
	Models      []ModelSpec `yaml:"models"`
	Routes      []RouteSpec `yaml:"routes"`
}
