package layout

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
	Integration string `yaml:"integration"`
	Deployment  string `yaml:"deployment"`
}

// Walk implements post-node visiting of each element in the manifest.
//
// Every App, Service and Environment is called once, and any error from the
// handling function terminates the Walk.
func (m Manifest) Walk(visitor ManifestVisitor) error {
	for _, env := range m.Environments {
		for _, app := range env.Apps {
			for _, svc := range app.Services {
				err := visitor.Service(env, app, svc)
				if err != nil {
					return err
				}
			}
			err := visitor.Application(env, app)
			if err != nil {
				return err
			}
		}
		err := visitor.Environment(env)
		if err != nil {
			return err
		}
	}
	return nil
}
