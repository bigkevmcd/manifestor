package layout

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestManifestWalk(t *testing.T) {
	m := &Manifest{
		Environments: []*Environment{
			&Environment{
				Name: "development",
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
			&Environment{
				Name: "staging",
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

	v := &testVisitor{paths: []string{}}
	m.Walk(v)
	sort.Strings(v.paths)

	want := []string{
		"development/my-app-1",
		"development/my-app-1/app-1-service-http",
		"development/my-app-1/app-1-service-test",
		"development/my-app-2",
		"development/my-app-2/app-2-service",
		"envs/development",
		"envs/staging",
		"staging/my-app-1",
		"staging/my-app-1/app-1-service-user",
	}

	if diff := cmp.Diff(want, v.paths); diff != "" {
		t.Fatalf("tree files: %s", diff)
	}
}

type testVisitor struct {
	paths []string
}

func (v *testVisitor) Service(env *Environment, app *Application, svc *Service) error {
	v.paths = append(v.paths, filepath.Join(env.Name, app.Name, svc.Name))
	return nil
}

func (v *testVisitor) Application(env *Environment, app *Application) error {
	v.paths = append(v.paths, filepath.Join(env.Name, app.Name))
	return nil
}

func (v *testVisitor) Environment(env *Environment) error {
	v.paths = append(v.paths, filepath.Join("envs", env.Name))
	return nil
}
