# Manifestor specification

This is the current specification for Manifestor manifests.

```yaml
# top-level "environments" a list of environments, probably representing
# namespaces, possibly clusters.
environments:
  # this is an environment called "development"
  - name: development
    # "pipelines" is a property of the environment, this represents the
    # pipeline bindings to be used in the automatically generated EventListener.
    # If a service has a `source_url` then we can automatically trigger these
    # pipeline bindings based on Git Hosting Service hooks.
    # 
    # The "integration" pipeline would be triggered on PullRequest events.
    # The "deployment" pipeline would be triggered on Push events.
    pipelines:
      integration:
        template: dev-ci-template
        binding: dev-ci-binding
      deployment:
        template: dev-cd-template
        binding: dev-cd-binding
    # This is a set of "apps", which according to the IBM/Red Hat repository
    # "contain" a set of services.
    apps:
      - name: my-app-1
        services:
          - name: app-1-service-http
            # the "source_url" represents the "upstream" source code e.g. the
            # Java/Go/Node.js code.
            #
            # This is used when generating the EventListener to drive the CI
            # process, PullRequests and Pushes to this branch drive the
            # pipelines.
            source_url: https://github.com/myproject/myservice.git
          - name: app-1-service-metrics
      - name: my-app-2
        # this is the set of "services" associated with this app "my-app-1".
        services:
          - name: app-2-service
  - name: staging
    apps:
      - name: my-app-1
        services:
          - name: app-1-service-http
  - name: production
    apps:
      - name: my-app-1
        services:
          - name: app-1-service-http
          - name: app-1-service-metrics
            # the "config_repo" property is an upstream reference to
            # configuration in another repository.
            #
            # This is used by the ArgoCD creation to drive deployments of
            # configuration from this repository.
            config_repo:
              url: https://github.com/testing/testing
              ref: master
              path: config
```

## Pipelines

Pipelines could be defined at a top level with a name, and binding/template and
referenced in the services by name, rather than by template.

Pipelines could also be specified at the individual services level, so for
example, you could specify a pipeline for CI on a specific service.
