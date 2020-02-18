package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	divingbell "github.com/tdaines42/diving-bell/pkg/diving-bell"
)

var write bool

func init() {
	configGenCmd.Flags().BoolVar(&write, "write", false, "Write the config to the local file")
	configCmd.AddCommand(configGenCmd)
}

var configGenCmd = &cobra.Command{
	Use:   "gen [cluster name] [terraform workspace path]",
	Short: "Generate and print the config to the console",
	Long:  ``,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if write {
			divingbell.UpdateClusterConfigFile(args[0], args[1])
		} else {
			fmt.Println(divingbell.ClusterConfigYamlString(args[0], args[1]))
		}
	},
}
