package layout

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

func TestManifestPaths(t *testing.T) {
	m := &manifest.Manifest{
		Environments: map[string]*manifest.Environment{
			"development": &manifest.Environment{
				Apps: []*manifest.Application{
					&manifest.Application{
						Name: "my-app-1",
						Services: []*manifest.Service{
							&manifest.Service{Name: "app-1-service-http"},
							&manifest.Service{Name: "app-1-service-test"},
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
							&manifest.Service{Name: "app-1-service-user"},
						},
					},
				},
			},
		},
	}

	paths := manifestPaths(m)
	sort.Strings(paths)
	want := []string{
		"development/apps/my-app-1/base/kustomization.yaml",
		"development/apps/my-app-1/overlays/kustomization.yaml",
		"development/apps/my-app-2/base/kustomization.yaml",
		"development/apps/my-app-2/overlays/kustomization.yaml",
		"development/services/app-1-service-http/base/config/kustomization.yaml",
		"development/services/app-1-service-http/base/kustomization.yaml",
		"development/services/app-1-service-http/overlays/kustomization.yaml",
		"development/services/app-1-service-test/base/config/kustomization.yaml",
		"development/services/app-1-service-test/base/kustomization.yaml",
		"development/services/app-1-service-test/overlays/kustomization.yaml",
		"development/services/app-2-service/base/config/kustomization.yaml",
		"development/services/app-2-service/base/kustomization.yaml",
		"development/services/app-2-service/overlays/kustomization.yaml",
		"staging/apps/my-app-1/base/kustomization.yaml",
		"staging/apps/my-app-1/overlays/kustomization.yaml",
		"staging/services/app-1-service-user/base/config/kustomization.yaml",
		"staging/services/app-1-service-user/base/kustomization.yaml",
		"staging/services/app-1-service-user/overlays/kustomization.yaml",
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
