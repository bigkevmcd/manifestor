package layout

// ManifestVisitor is an interface for accessing all components of a manifest in
// turn.
type ManifestVisitor interface {
	Environment(*Environment) error
	Application(*Environment, *Application) error
	Service(*Environment, *Application, *Service) error
}
