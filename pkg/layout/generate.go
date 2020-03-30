package layout

import (
	"fmt"
	"path/filepath"

	"github.com/bigkevmcd/manifestor/pkg/manifest"
)

func filesForService() map[string]interface{} {
	return map[string]interface{}{
		"base/kustomization.yaml":        serviceBaseKustomization(),
		"base/config/kustomization.yaml": serviceConfigKustomization(),
		"overlays/kustomization.yaml":    serviceOverlaysKustomization(),
	}
}

func serviceBaseKustomization() map[string]interface{} {
	return map[string]interface{}{
		"bases": []string{"./config"},
	}
}

func serviceConfigKustomization() map[string]interface{} {
	return map[string]interface{}{
		"resources": []string{""},
	}
}

func serviceOverlaysKustomization() map[string]interface{} {
	return map[string]interface{}{
		"bases": []string{"../config"},
	}
}

func appKustomization(services []string) map[string]interface{} {
	overlayPaths := make([]string, len(services))
	for i, s := range services {
		overlayPaths[i] = fmt.Sprintf("../../../services/%s/overlays", s)
	}
	return map[string]interface{}{
		"base/kustomization.yaml": map[string]interface{}{
			"bases": overlayPaths,
		},
		"overlays/kustomization.yaml": []string{"../config"},
	}
}

func environmentFiles(apps []string) map[string]interface{} {
	overlayPaths := make([]string, len(apps))
	for i, a := range apps {
		overlayPaths[i] = fmt.Sprintf("../../../apps/%s/overlays", a)
	}
	return map[string]interface{}{
		"base/kustomization.yaml": map[string]interface{}{
			"bases": overlayPaths,
		},
		"overlays/kustomization.yaml": []string{"../base"},
	}

}

func manifestPaths(man *manifest.Manifest) []string {
	files := []string{}
	appNames := []string{}
	for name, env := range man.Environments {
		for _, app := range env.Apps {
			serviceNames := []string{}
			for _, svc := range app.Services {
				servicePath := filepath.Join(name, "services", svc.Name)
				for f, _ := range filesForService() {
					files = append(files, filepath.Join(servicePath, f))
				}
			}
			for k, _ := range appKustomization(serviceNames) {
				files = append(files, filepath.Join(name, "apps", app.Name, k))
			}
			appNames = append(appNames, app.Name)
		}
	}
	for k, _ := range environmentFiles(appNames) {
		files = append(files, filepath.Join("env", k))
	}
	return files
}
