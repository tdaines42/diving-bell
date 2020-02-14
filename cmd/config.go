package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	divingbell "github.com/tdaines42/diving-bell/pkg/diving-bell"
)

var write bool
var retrieve bool
var update bool

func init() {
	configCmd.Flags().BoolVar(&write, "write", false, "Write the config to the local file")
	configCmd.Flags().BoolVar(&retrieve, "retrieve", false, "Retrieve the config from the cluster")
	configCmd.Flags().BoolVar(&update, "update", false, "Update the config in the cluster")
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config [cluster name] [terraform workspace path]",
	Short: "Print the config to the console",
	Long:  ``,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if write {
			divingbell.UpdateClusterConfigFile(args[0], args[1])
		} else if retrieve {
			divingbell.RetrieveClusterConfig(args[0])
		} else if update {
			divingbell.StoreClusterConfig(args[0], args[1])
		} else {
			fmt.Println(divingbell.ClusterConfigYamlString(args[0], args[1]))
		}
	},
}
