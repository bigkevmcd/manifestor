package layout

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestManifestWalk(t *testing.T) {
	m := &Manifest{
		Environments: map[string]*Environment{
			"development": &Environment{
				Apps: []*Application{
					&Application{
						Name: "my-app-1",
						Services: []*Service{
							&Service{Name: "app-1-service-http"},
							&Service{Name: "app-1-service-test"},
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
							&Service{Name: "app-1-service-user"},
						},
					},
				},
			},
		},
	}

	paths := []string{}
	m.Walk(func(env string, app *Application, service *Service) error {
		paths = append(paths, filepath.Join(env, app.Name, service.Name))
		return nil
	})
	sort.Strings(paths)

	want := []string{
		"development/my-app-1/app-1-service-http",
		"development/my-app-1/app-1-service-test",
		"development/my-app-2/app-2-service",
		"staging/my-app-1/app-1-service-user",
	}

	if diff := cmp.Diff(want, paths); diff != "" {
		t.Fatalf("tree files: %s", diff)
	}
}
