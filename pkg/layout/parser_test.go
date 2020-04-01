package layout

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	parseTests := []struct {
		filename string
		want     *Manifest
	}{
		{"testdata/example1.yaml", &Manifest{
			Environments: []*Environment{
				&Environment{
					Name: "development",
					Apps: []*Application{
						&Application{
							Name: "my-app-1",
							Services: []*Service{
								&Service{Name: "app-1-service-http"},
								&Service{Name: "app-1-service-metrics"},
							},
						},
						&Application{
							Name: "my-app-2",
							Services: []*Service{
								&Service{Name: "app-2-service"},
							},
						},
					},
				},
				&Environment{
					Name: "staging",
					Apps: []*Application{
						&Application{Name: "my-app-1",
							Services: []*Service{
								&Service{Name: "app-1-service-http"},
							},
						},
					},
				},
				&Environment{
					Name: "production",
					Apps: []*Application{
						&Application{Name: "my-app-1",
							Services: []*Service{
								&Service{Name: "app-1-service-http"},
								&Service{Name: "app-1-service-metrics", Repo: &Repository{
									SourceURL: "https://github.com/testing/testing",
									Ref:       "master",
									Path:      "config",
								},
								},
							},
						},
					},
				},
			},
		},
		},
	}

	for _, tt := range parseTests {
		t.Run(fmt.Sprintf("parsing %s", tt.filename), func(rt *testing.T) {
			f, err := os.Open(tt.filename)
			if err != nil {
				rt.Fatalf("failed to open %v: %s", tt.filename, err)
			}
			defer f.Close()

			got, err := Parse(f)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				rt.Errorf("Parse(%s) failed diff\n%s", tt.filename, diff)
			}
		})
	}
}
