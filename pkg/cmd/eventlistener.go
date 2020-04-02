package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/bigkevmcd/manifestor/pkg/eventlistener"
	"github.com/bigkevmcd/manifestor/pkg/layout"
)

func makeEventListenerCommand() *cobra.Command {
	var eventListenerName string
	cmd := &cobra.Command{
		Use:   "eventlistener <manifest file>",
		Short: "generate an eventlistener",
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

			el := eventlistener.GenerateEventListener(args[1], m)
			b, err := yaml.Marshal(el)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", b)
			return nil
		},
	}
	cmd.Flags().StringVar(&eventListenerName, "eventlistener-name", "default-event-listener", "provide a name for the generate Tekton Triggers EventListener")
	return cmd
}
