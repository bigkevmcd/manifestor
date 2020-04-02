package argoapps

import (
	"testing"

	argoappv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bigkevmcd/manifestor/pkg/layout"
)

func TestMakeApplication(t *testing.T) {
	got := makeApplication("blue-app", "default", "testing-ns", "https://kubernetes.default.svc", argoappv1.ApplicationSource{RepoURL: "https://example.com/myorg/myproject.git", Path: "services/service-1/base"})

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

func TestGenerateArgoApplications(t *testing.T) {
	testManifest := mustParseManifest(t, "../testdata/example2.yaml")
	repoURL := "https://example.com/myorg/myrepo.git"
	apps, err := GenerateArgoApplications(repoURL, testManifest)
	if err != nil {
		t.Fatal(err)
	}
	want := []*argoappv1.Application{
		&argoappv1.Application{
			TypeMeta: applicationTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-app-1-service-http",
			},
			Spec: argoappv1.ApplicationSpec{
				Project: "default",
				Destination: argoappv1.ApplicationDestination{
					Namespace: "dev",
					Server:    "https://kubernetes.default.svc",
				},
				Source: argoappv1.ApplicationSource{
					RepoURL: repoURL,
					Path:    "dev/services/service-http/base",
				},
				SyncPolicy: syncPolicy,
			},
		},
		&argoappv1.Application{
			TypeMeta: applicationTypeMeta,
			ObjectMeta: metav1.ObjectMeta{
				Name: "my-app-2-service-metrics",
			},
			Spec: argoappv1.ApplicationSpec{
				Project: "default",
				Destination: argoappv1.ApplicationDestination{
					Namespace: "production",
					Server:    "https://kubernetes.default.svc",
				},
				Source: argoappv1.ApplicationSource{
					RepoURL:        "https://example.com/testing/testing.git",
					Path:           "config",
					TargetRevision: "master",
				},
				SyncPolicy: syncPolicy,
			},
		},
	}
	if diff := cmp.Diff(want, apps); diff != "" {
		t.Errorf("generate failed diff\n%s", diff)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func mustParseManifest(t *testing.T, filename string) *layout.Manifest {
	t.Helper()
	m, err := layout.ParseFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return m
}
