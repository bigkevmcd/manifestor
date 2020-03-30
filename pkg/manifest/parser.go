package manifest

import (
	"io"

	"gopkg.in/yaml.v2"
)

type Repository struct {
	SourceURL string `yaml:"url"`
	Ref       string `yaml:"ref"`
	Path      string `yaml:"path"`
}

type Service struct {
	Name string      `yaml:"name"`
	Repo *Repository `yaml:"repo"`
}

type Application struct {
	Name     string     `yaml:"name"`
	Services []*Service `yaml:"services"`
}

type Environment struct {
	Apps []*Application `yaml:"apps"`
}
type Manifest struct {
	Environments map[string]*Environment `yaml:"environments"`
}

// Parse decodes YAML describing an environment manifest.
func Parse(in io.Reader) (*Manifest, error) {
	dec := yaml.NewDecoder(in)
	m := &Manifest{}
	err := dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
