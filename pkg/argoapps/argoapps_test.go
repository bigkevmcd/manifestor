package argoapps

import (
	"testing"

	argoappv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bigkevmcd/manifestor/pkg/layout"
)

var (
	testService = &layout.Service{
		Name:      "app-1-service-http",
		SourceURL: "https://github.com/testing/testing",
	}

	testApplication = &layout.Application{
		Name: "my-app-1",
		Services: []*layout.Service{
			testService,
		},
	}

	testEnvironment = &layout.Environment{
		Name: "development",
		Apps: []*layout.Application{
			testApplication,
		},
	}

	testManifest = &layout.Manifest{
		Environments: []*layout.Environment{
			testEnvironment,
		},
	}
)

func TestMakeApplication(t *testing.T) {
	got := makeApplication("blue-app", "default", "testing-ns", "https://kubernetes.default.svc", "https://example.com/myorg/myproject.git", "services/service-1/base")

	want := &argoappv1.Application{
		TypeMeta: applicationTypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name: "blue-app",
		},
		Spec: argoappv1.ApplicationSpec{
			Project: "default",
			Destination: argoappv1.ApplicationDestination{
				Namespace: "testing-ns",
				Server:    "https://kubernetes.default.svc",
			},
			Source: argoappv1.ApplicationSource{
				RepoURL: "https://example.com/myorg/myproject.git",
				Path:    "services/service-1/base",
			},
			SyncPolicy: syncPolicy,
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("makeApplication() diff: %s\n", diff)
	}
}

func TestExtractServices(t *testing.T) {
	want := []service{
		{"my-app-1-app-1-service-http", layout.PathForService(testEnvironment, testApplication, testService), "staging"},
	}
	l, err := extractServices(testManifest)
	assertNoError(t, err)

	if diff := cmp.Diff(want, l); diff != "" {
		t.Errorf("generate failed diff\n%s", diff)
	}
}

func TestGenerateArgoApplications(t *testing.T) {
	el := GenerateArgoApplications("https://example.com/myorg/myrepo.git", testManifest)
	want := []*argoappv1.Application{}

	if diff := cmp.Diff(want, el); diff != "" {
		t.Errorf("generate failed diff\n%s", diff)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
