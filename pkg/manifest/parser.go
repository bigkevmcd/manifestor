package manifest

import (
	"io"

	"gopkg.in/yaml.v2"
)

type Service struct {
	Name string `json:"name"`
}

type Application struct {
	Name     string     `json:"name"`
	Services []*Service `json:"services"`
}

type Environment struct {
	Apps []*Application `json:"apps"`
}
type Manifest struct {
	Environments map[string]*Environment `json:"environments"`
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
