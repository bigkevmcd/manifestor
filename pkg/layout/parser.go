package layout

import (
	"io"

	"gopkg.in/yaml.v2"
)

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
