package layout_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/bigkevmcd/manifestor/pkg/manifest"
	"github.com/google/go-cmp/cmp"
)

func manifestPaths(man *manifest.Manifest) []string {
	files := []string{}
	for name, env := range man.Environments {
		for _, app := range env.Apps {
			for _, svc := range app.Services {
				files = append(files, filepath.Join(name, app.Name, svc.Name))
			}
		}
	}
	return files
}

func TestGenerate(t *testing.T) {
	m := &manifest.Manifest{
		Environments: map[string]*manifest.Environment{
			"development": &manifest.Environment{
				Apps: []*manifest.Application{
					&manifest.Application{
						Name: "my-app-1",
						Services: []*manifest.Service{
							&manifest.Service{Name: "app-1-service-http"},
						},
					},
					&manifest.Application{
						Name: "my-app-2",
						Services: []*manifest.Service{
							&manifest.Service{Name: "app-2-service"},
						},
					},
				},
			},
			"staging": &manifest.Environment{
				Apps: []*manifest.Application{
					&manifest.Application{Name: "my-app-1",
						Services: []*manifest.Service{
							&manifest.Service{Name: "app-1-service-http"},
						},
					},
				},
			},
		},
	}

	paths := manifestPaths(m)
	sort.Strings(paths)
	want := []string{
		"development/my-app-1/app-1-service-http",
		"development/my-app-2/app-2-service",
		"staging/my-app-1/app-1-service-http",
	}
	if diff := cmp.Diff(want, paths); diff != "" {
		t.Fatalf("tree files: %s", diff)
	}
}

func tempDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := ioutil.TempDir(os.TempDir(), "gnome")
	assertNoError(t, err)
	return dir, func() {
		err := os.RemoveAll(dir)
		assertNoError(t, err)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertTreeFiles(t *testing.T, path string, want []string) {
	t.Helper()
	got := []string{}
	filepath.Walk(path, func(treepath string, info os.FileInfo, err error) error {
		relativePath := strings.TrimPrefix(treepath, path)
		if relativePath != "" {
			got = append(got, relativePath)
		}
		return err
	})

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("tree files: %s", diff)
	}
}
