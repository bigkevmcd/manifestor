module github.com/bigkevmcd/manifestor

go 1.14

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/argoproj/argo-cd v1.5.0-rc3
	github.com/argoproj/pkg v0.0.0-20200319004004-f46beff7cd54 // indirect
	github.com/google/go-cmp v0.4.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/robfig/cron v1.2.0 // indirect
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	github.com/tektoncd/pipeline v0.10.1 // indirect
	github.com/tektoncd/triggers v0.3.1
	gomodules.xyz/jsonpatch/v2 v2.1.0 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.0 // indirect
	k8s.io/apimachinery v0.18.0
	knative.dev/pkg v0.0.0-20200331145051-94b2b6aaaf4b // indirect
	sigs.k8s.io/yaml v1.2.0
)

// Pin k8s deps to 1.16.5
replace (
	k8s.io/api => k8s.io/api v0.16.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.5
)

replace k8s.io/client-go => k8s.io/client-go v0.16.5
