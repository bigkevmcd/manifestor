package eventlistener

import (
	"testing"

	"github.com/bigkevmcd/manifestor/pkg/layout"
	"github.com/google/go-cmp/cmp"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var testManifest = &layout.Manifest{
	Environments: []*layout.Environment{
		&layout.Environment{
			Name: "development",
			Pipelines: &layout.Pipelines{
				Integration: "test-ci",
				Deployment:  "test-cd",
			},
			Apps: []*layout.Application{
				&layout.Application{
					Name: "my-app-1",
					Services: []*layout.Service{
						&layout.Service{
							Name:      "app-1-service-http",
							SourceURL: "https://github.com/testing/testing",
						},
					},
				},
			},
		},
	},
}

func TestParseManifest(t *testing.T) {
	want := []service{{"testing/testing", "my-app-1-app-1-service-http", "development", "test-ci", "test-cd"}}
	l, err := parseManifest(testManifest)
	assertNoError(t, err)

	if diff := cmp.Diff(want, l); diff != "" {
		t.Errorf("generate failed diff\n%s", diff)
	}
}

func TestGenerateEventListener(t *testing.T) {
	el := GenerateEventListener("my-test", testManifest)
	want := &triggersv1.EventListener{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EventListener",
			APIVersion: "triggers.tekton.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-test",
		},
		Spec: triggersv1.EventListenerSpec{
			Triggers: []triggersv1.EventListenerTrigger{
				triggersv1.EventListenerTrigger{
					Bindings: []*triggersv1.EventListenerBinding{
						&triggersv1.EventListenerBinding{
							Name: "test-ci",
						},
					},
					Template: triggersv1.EventListenerTemplate{Name: "test-cd"},
					Name:     "my-app-1-app-1-service-http",
					Interceptors: []*triggersv1.EventInterceptor{
						&triggersv1.EventInterceptor{
							CEL: &triggersv1.CELInterceptor{
								Filter: "(header.match('X-GitHub-Event', 'pull_request') && body.action == 'opened' || body.action == 'synchronize') && body.pull_request.head.repo.full_name == 'testing/testing'",
							},
						},
					},
				},
			},
		},
	}

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
