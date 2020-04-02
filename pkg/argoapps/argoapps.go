package argoapps

import (
	"fmt"

	argoappv1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bigkevmcd/manifestor/pkg/layout"
)

var (
	applicationTypeMeta = metav1.TypeMeta{
		Kind:       "Application",
		APIVersion: "argoproj.io/v1alpha1",
	}

	syncPolicy = &argoappv1.SyncPolicy{
		Automated: &argoappv1.SyncPolicyAutomated{
			Prune:    true,
			SelfHeal: true,
		},
	}
)

const defaultServer = "https://kubernetes.default.svc"

func GenerateArgoApplications(repoURL string, m *layout.Manifest) []*argoappv1.Application {
	l, _ := extractServices(m)
	apps := make([]*argoappv1.Application, len(l))
	for i, s := range l {
		apps[i] = makeApplication(s.AppName, "default", s.Environment, defaultServer, repoURL, s.Path)
	}
	return apps
}

// TODO: structure these arguments a bit better.
func makeApplication(appName, project, ns, server, repoURL, path string) *argoappv1.Application {
	return &argoappv1.Application{
		TypeMeta: applicationTypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name: appName,
		},
		Spec: argoappv1.ApplicationSpec{
			Project: project,
			Destination: argoappv1.ApplicationDestination{
				Namespace: ns,
				Server:    server,
			},
			Source: argoappv1.ApplicationSource{
				RepoURL: repoURL,
				Path:    path,
			},
			SyncPolicy: syncPolicy,
		},
	}
}

type serviceVisitor struct {
	services []service
}

func (ev *serviceVisitor) Service(env *layout.Environment, app *layout.Application, svc *layout.Service) error {
	ev.services = append(ev.services, service{fmt.Sprintf("%s-%s", app.Name, svc.Name), layout.PathForService(env, app, svc), env.Name})
	return nil
}

func (ev *serviceVisitor) Application(env *layout.Environment, app *layout.Application) error {
	return nil
}

func (ev *serviceVisitor) Environment(env *layout.Environment) error {
	return nil
}

type service struct {
	AppName     string
	Path        string
	Environment string
}

func extractServices(m *layout.Manifest) ([]service, error) {
	l := &serviceVisitor{services: []service{}}
	m.Walk(l)
	return l.services, nil
}
