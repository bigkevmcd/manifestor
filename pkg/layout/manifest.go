package layout

import (
	"sort"
)

// Manifest describes a set of environments, apps and services for deployment.
type Manifest struct {
	Environments []*Environment `yaml:"environments"`
}

// Environment is a slice of Apps, these are the named apps in the namespace.
type Environment struct {
	Name      string         `yaml:"name"`
	Pipelines *Pipelines     `yaml:"pipelines"`
	Apps      []*Application `yaml:"apps"`
}

func (e Environment) GoString() string {
	return e.Name
}

// Application has many services.
type Application struct {
	Name     string     `yaml:"name"`
	Services []*Service `yaml:"services"`
}

// Service has an upstream source.
type Service struct {
	Name       string      `yaml:"name"`
	SourceURL  string      `yaml:"source_url"`
	ConfigRepo *Repository `yaml:"config_repo"`
}

// Repository refers to an upstream source for reading additional config from.
type Repository struct {
	URL  string `yaml:"url"`
	Ref  string `yaml:"ref"`
	Path string `yaml:"path"`
}

// Pipelines describes the names for pipelines to be executed for CI and CD.
//
// These pipelines will be executed with a Git clone URL and commit SHA.
type Pipelines struct {
	Integration *TemplateBinding `yaml:"integration"`
	Deployment  *TemplateBinding `yaml:"deployment"`
}

// TemplateBinding is a combination of the template and binding to be used for a
// pipeline execution.
type TemplateBinding struct {
	Template string `yaml:"template"`
	Binding  string `yaml:"binding"`
}

// Walk implements post-node visiting of each element in the manifest.
//
// Every App, Service and Environment is called once, and any error from the
// handling function terminates the Walk.
//
// The environments are sorted using a custom sorting mechanism, that orders by
// name, but, moves CICD environments to the bottom of the list.
func (m Manifest) Walk(visitor interface{}) error {
	sort.Sort(byName(m.Environments))
	for _, env := range m.Environments {
		for _, app := range env.Apps {
			for _, svc := range app.Services {
				if v, ok := visitor.(ServiceVisitor); ok {
					err := v.Service(env, app, svc)
					if err != nil {
						return err
					}
				}
			}
			if v, ok := visitor.(ApplicationVisitor); ok {
				err := v.Application(env, app)
				if err != nil {
					return err
				}
			}
		}
		if v, ok := visitor.(EnvironmentVisitor); ok {
			err := v.Environment(env)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

type byName []*Environment

func (a byName) Len() int      { return len(a) }
func (a byName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool {
	if a[j].Name == "cicd" {
		return true
	}
	return a[i].Name < a[j].Name
}
