package eventlistener

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bigkevmcd/manifestor/pkg/layout"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ciFilter = "(header.match('X-GitHub-Event', 'pull_request') && body.action == 'opened' || body.action == 'synchronize') && body.pull_request.head.repo.full_name == '%s'"

func GenerateEventListener(n string, m *layout.Manifest) *triggersv1.EventListener {
	return &triggersv1.EventListener{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EventListener",
			APIVersion: "triggers.tekton.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: n,
		},
		Spec: triggersv1.EventListenerSpec{
			Triggers: makeEventListenerTriggers(m),
		},
	}
}

// TODO: deal with the error here.
func makeEventListenerTriggers(m *layout.Manifest) []triggersv1.EventListenerTrigger {
	l, _ := extractServices(m)
	triggers := make([]triggersv1.EventListenerTrigger, len(l))
	for i, s := range l {
		triggers[i] = createListenerTrigger(s.Name, ciFilter, s.RepoName, s.CI.Binding, s.CI.Template)
	}
	return triggers
}

type serviceVisitor struct {
	services []service
}

func extractRepo(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	parts := strings.Split(parsed.Path, "/")
	return fmt.Sprintf("%s/%s", parts[1], strings.TrimSuffix(parts[2], ".git")), nil
}

func (ev *serviceVisitor) Service(env *layout.Environment, app *layout.Application, svc *layout.Service) error {
	if svc.SourceURL == "" {
		return nil
	}
	repo, err := extractRepo(svc.SourceURL)
	if err != nil {
		return err
	}
	ev.services = append(ev.services, service{repo, fmt.Sprintf("%s-%s", app.Name, svc.Name), env.Name, env.Pipelines.Integration, env.Pipelines.Deployment})
	return nil
}

type service struct {
	RepoName string
	Name     string
	Env      string
	CI       *layout.TemplateBinding
	CD       *layout.TemplateBinding
}

func extractServices(m *layout.Manifest) ([]service, error) {
	l := &serviceVisitor{services: []service{}}
	m.Walk(l)
	return l.services, nil
}
