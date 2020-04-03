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

func TestManifestWalkCallsCICDEnvironmentLast(t *testing.T) {
	m := &Manifest{
		Environments: []*Environment{
			&Environment{
				Name: "cicd",
			},
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

	want := []string{
		"development/my-app-1/app-1-service-http",
		"development/my-app-1/app-1-service-test",
		"development/my-app-1",
		"development/my-app-2/app-2-service",
		"development/my-app-2",
		"envs/development",
		"staging/my-app-1/app-1-service-user",
		"staging/my-app-1",
		"envs/staging",
		"cicd/development/app-1-service-http",
		"cicd/development/app-1-service-test",
		"cicd/development/app-2-service",
		"cicd/staging/app-1-service-user",
		"envs/cicd",
	}

	if diff := cmp.Diff(want, v.paths); diff != "" {
		t.Fatalf("tree files: %s", diff)
	}
}

func TestEnviromentSorting(t *testing.T) {
	envNames := func(envs []*Environment) []string {
		n := make([]string, len(envs))
		for i, v := range envs {
			n[i] = v.Name
		}
		return n
	}
	makeEnvs := func(ns []string) []*Environment {
		n := make([]*Environment, len(ns))
		for i, v := range ns {
			n[i] = &Environment{Name: v}
		}
		return n

	}
	envTests := []struct {
		names []string
		want  []string
	}{
		{[]string{"prod", "staging", "dev"}, []string{"dev", "prod", "staging"}},
		{[]string{"cicd", "staging", "dev"}, []string{"dev", "staging", "cicd"}},
	}

	for _, tt := range envTests {
		envs := makeEnvs(tt.names)
		sort.Sort(byName(envs))
		if diff := cmp.Diff(envNames(envs), tt.want); diff != "" {
			t.Errorf("sort(%#v): %s", envs, diff)
		}
	}
}

type testVisitor struct {
	pipelineServices []string
	paths            []string
}

func (v *testVisitor) Service(env *Environment, app *Application, svc *Service) error {
	v.paths = append(v.paths, filepath.Join(env.Name, app.Name, svc.Name))
	v.pipelineServices = append(v.pipelineServices, filepath.Join("cicd", env.Name, svc.Name))
	return nil
}

func (v *testVisitor) Application(env *Environment, app *Application) error {
	v.paths = append(v.paths, filepath.Join(env.Name, app.Name))
	return nil
}

func (v *testVisitor) Environment(env *Environment) error {
	if env.Name == "cicd" {
		v.paths = append(v.paths, v.pipelineServices...)
	}
	v.paths = append(v.paths, filepath.Join("envs", env.Name))
	return nil
}
