package argoapps

import (
	"fmt"
	"path/filepath"

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

const (
	defaultServer  = "https://kubernetes.default.svc"
	defaultProject = "default"
)

func GenerateArgoApplications(repoURL string, m *layout.Manifest) ([]*argoappv1.Application, error) {
	l := &visitor{applications: []*argoappv1.Application{}, repoURL: repoURL}
	err := m.Walk(l)
	if err != nil {
		return nil, err
	}
	return l.applications, nil
}

// TODO: structure these arguments a bit better.
func makeApplication(appName, project, ns, server string, source argoappv1.ApplicationSource) *argoappv1.Application {
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
			Source:     source,
			SyncPolicy: syncPolicy,
		},
	}
}

type visitor struct {
	repoURL      string
	applications []*argoappv1.Application
}

func (sv *visitor) Service(env *layout.Environment, app *layout.Application, svc *layout.Service) error {
	newApp := makeApplication(
		fmt.Sprintf("%s-%s", app.Name, svc.Name),
		defaultProject,
		env.Name,
		defaultServer,
		sv.makeSource(env, app, svc),
	)
	sv.applications = append(sv.applications, newApp)
	return nil
}

func (sv *visitor) makeSource(env *layout.Environment, app *layout.Application, svc *layout.Service) argoappv1.ApplicationSource {
	if svc.ConfigRepo == nil {
		return argoappv1.ApplicationSource{
			RepoURL: sv.repoURL,
			Path:    filepath.Join(layout.PathForService(env, app, svc), "base"),
		}
	}
	return argoappv1.ApplicationSource{
		RepoURL:        svc.ConfigRepo.URL,
		Path:           svc.ConfigRepo.Path,
		TargetRevision: svc.ConfigRepo.Ref,
	}
}
