package manifest

import (
	"io"

	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func TestParse(t *testing.T) {
	parseTests := []struct {
		filename string
		want     *Manifest
	}{
		{"testdata/example1.yaml", &Manifest{
			Environments: map[string]*Environment{
				"development": &Environment{
					Apps: []*Application{
						&Application{
							Name: "my-app-1",
							Services: []*Service{
								&Service{Name: "app-1-service-http"},
								&Service{Name: "app-1-service-metrics"},
							},
						},
						&Application{
							Name: "my-app-2",
							Services: []*Service{
								&Service{Name: "app-2-service"},
							},
						},
					},
				},
				"staging": &Environment{
					Apps: []*Application{
						&Application{Name: "my-app-1",
							Services: []*Service{
								&Service{Name: "app-1-service-http"},
							},
						},
					},
				},
				"production": &Environment{
					Apps: []*Application{
						&Application{Name: "my-app-1",
							Services: []*Service{
								&Service{Name: "app-1-service-http"},
								&Service{Name: "app-1-service-metrics"},
							},
						},
					},
				},
			},
		},
		},
	}

	for _, tt := range parseTests {
		t.Run(fmt.Sprintf("parsing %s", tt.filename), func(rt *testing.T) {
			f, err := os.Open(tt.filename)
			if err != nil {
				rt.Fatalf("failed to open %v: %s", tt.filename, err)
			}
			defer f.Close()

			got, err := Parse(f)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				rt.Errorf("Parse(%s) failed diff\n%s", tt.filename, diff)
			}
		})
	}
}
