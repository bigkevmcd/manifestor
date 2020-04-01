package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/bigkevmcd/manifestor/pkg/eventlistener"
	"github.com/bigkevmcd/manifestor/pkg/layout"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func makeEventListenerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "eventlistener",
		Short: "generate an eventlistener",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			f, err := os.Open(args[0])
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			m, err := layout.Parse(f)
			if err != nil {
				log.Fatal(err)
			}

			el := eventlistener.GenerateEventListener(args[1], m)
			b, err := yaml.Marshal(el)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", b)
		},
	}
}
