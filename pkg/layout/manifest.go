package layout

// Manifest describes a set of environments, apps and services for deployment.
type Manifest struct {
	Environments map[string]*Environment `yaml:"environments"`
}

// Environment is a slice of Apps, these are the named apps in the namespace.
type Environment struct {
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

type serviceVisitor func(string, *Application, *Service) error

func (m Manifest) Walk(visitor serviceVisitor) error {
	for envName, env := range m.Environments {
		for _, app := range env.Apps {
			for _, svc := range app.Services {
				err := visitor(envName, app, svc)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
