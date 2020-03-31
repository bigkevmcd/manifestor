package layout

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Bootstrap takes a manifest and a prefix, and writes the files from the manifest
// starting with the prefix.
func Bootstrap(prefix string, m *Manifest) error {
	appNames := []string{}
	for envName, env := range m.Environments {
		for _, app := range env.Apps {
			serviceNames := []string{}
			for _, svc := range app.Services {
				servicePath := filepath.Join(envName, "services", svc.Name)
				for f, v := range filesForService() {
					filename := filepath.Join(servicePath, f)
					err := writeWithPrefix(prefix, filename, v)
					if err != nil {
						return err
					}
				}
			}
			for k, v := range appKustomization(serviceNames) {
				filename := filepath.Join(envName, "apps", app.Name, k)
				err := writeWithPrefix(prefix, filename, v)
				if err != nil {
					return err
				}
			}
			appNames = append(appNames, app.Name)
		}
		for k, v := range environmentFiles(appNames) {
			filename := filepath.Join(envName, "env", k)
			err := writeWithPrefix(prefix, filename, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeWithPrefix(prefix, filename string, v interface{}) error {
	prefixedName := filepath.Join(prefix, filename)
	err := os.MkdirAll(filepath.Dir(prefixedName), 0755)
	if err != nil {
		return fmt.Errorf("failed to MkDirAll for %s: %v", filename, err)
	}
	f, err := os.Create(prefixedName)
	if err != nil {
		return fmt.Errorf("failed to Create file %s: %v", filename, err)
	}
	defer f.Close()
	return writeYAML(f, v)
}

func writeYAML(out io.Writer, v interface{}) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	_, err = fmt.Fprintf(out, "%s", data)
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}
	return nil
}

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
		"resources": []string{},
	}
}

func serviceOverlaysKustomization() map[string]interface{} {
	return map[string]interface{}{
		"bases": []string{"../base"},
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
		"overlays/kustomization.yaml": map[string]interface{}{
			"bases": []string{"../base"},
		},
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
