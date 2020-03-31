package layout

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBootstrap(t *testing.T) {
	tempDir, cleanup := makeTempDir(t)
	defer cleanup()
	m := &Manifest{
		Environments: map[string]*Environment{
			"development": &Environment{
				Apps: []*Application{
					&Application{
						Name: "my-app-1",
						Services: []*Service{
							&Service{Name: "app-1-service-http"},
						},
					},
				},
			},
		},
	}

	err := Bootstrap(tempDir, m)
	assertNoError(t, err)

	files := walkTree(t, tempDir)
	sort.Strings(files)
	want := []string{
		"/development",
		"/development/apps",
		"/development/apps/my-app-1",
		"/development/apps/my-app-1/base",
		"/development/apps/my-app-1/base/kustomization.yaml",
		"/development/apps/my-app-1/overlays",
		"/development/apps/my-app-1/overlays/kustomization.yaml",
		"/development/env",
		"/development/env/base",
		"/development/env/base/kustomization.yaml",
		"/development/env/overlays",
		"/development/env/overlays/kustomization.yaml",
		"/development/services",
		"/development/services/app-1-service-http",
		"/development/services/app-1-service-http/base",
		"/development/services/app-1-service-http/base/config",
		"/development/services/app-1-service-http/base/config/kustomization.yaml",
		"/development/services/app-1-service-http/base/kustomization.yaml",
		"/development/services/app-1-service-http/overlays",
		"/development/services/app-1-service-http/overlays/kustomization.yaml",
	}

	if diff := cmp.Diff(want, files); diff != "" {
		t.Errorf("bootstrap failed diff\n%s", diff)
	}
}

func makeTempDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := ioutil.TempDir(os.TempDir(), "manifest")
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

func walkTree(t *testing.T, start string) []string {
	walked := []string{}
	filepath.Walk(start, func(path string, info os.FileInfo, err error) error {
		trimmed := strings.TrimPrefix(path, start)
		if trimmed != "" {
			walked = append(walked, trimmed)
		}
		return err
	})
	return walked
}
