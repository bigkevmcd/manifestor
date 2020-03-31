package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bigkevmcd/manifestor/pkg/layout"
)

func init() {
	cobra.OnInitialize(initConfig)
}

var rootCmd = &cobra.Command{
	Use:   "manifestor",
	Short: "update a repository based on a manifest",
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

		err = layout.Bootstrap(args[1], m)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func initConfig() {
	viper.AutomaticEnv()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
