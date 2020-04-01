package layout

// Manifest describes a set of environments, apps and services for deployment.
type Manifest struct {
	Environments []*Environment `yaml:"environments"`
}

// Environment is a slice of Apps, these are the named apps in the namespace.
type Environment struct {
	Name string         `yaml:"name"`
	Apps []*Application `yaml:"apps"`
}

type Application struct {
	Name     string     `yaml:"name"`
	Services []*Service `yaml:"services"`
}

type Service struct {
	Name string      `yaml:"name"`
	Repo *Repository `yaml:"repo"`
}

type Repository struct {
	SourceURL string `yaml:"url"`
	Ref       string `yaml:"ref"`
	Path      string `yaml:"path"`
}

func (m Manifest) Walk(visitor ManifestVisitor) error {
	for _, env := range m.Environments {
		err := visitor.Environment(env)
		if err != nil {
			return err
		}
		for _, app := range env.Apps {
			err := visitor.Application(env, app)
			if err != nil {
				return err
			}
			for _, svc := range app.Services {
				err := visitor.Service(env, app, svc)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
