package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bigkevmcd/manifestor/pkg/layout"
)

func makeBootstrapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "bootstrap <manifest file> <output directory>",
		Short:        "write repository layout",
		Args:         cobra.MinimumNArgs(2),
		SilenceUsage: true,
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

			return layout.Bootstrap(args[1], m)
		},
	}
	return cmd
}
