package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/bigkevmcd/manifestor/pkg/argoapps"
	"github.com/bigkevmcd/manifestor/pkg/layout"
)

func makeArgoCDCommand() *cobra.Command {
	var repositoryURL string
	cmd := &cobra.Command{
		Use:   "argoapps <manifest file> --repository-url https://github.com/myorg/myrepository.git",
		Short: "generate ArgoCD application resources",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()

			m, err := layout.Parse(f)
			if err != nil {
				return err
			}

			el, err := argoapps.GenerateArgoApplications(repositoryURL, m)
			if err != nil {
				return err
			}
			b, err := yaml.Marshal(el)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", b)
			return nil
		},
	}
	cmd.Flags().StringVar(&repositoryURL, "repository-url", "", "full Git repository URL to deploy from")
	cmd.MarkFlagRequired("repository-url")
	return cmd
}
