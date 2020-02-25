package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	divingbell "github.com/tdaines42/diving-bell/pkg/diving-bell"
)

var write bool

func init() {
	configGenCmd.Flags().BoolVar(&write, "write", false, "Write the config to the local file")
	configGenCmd.Flags().StringVar(&kubernetesVersion, "kubernetes-version", "", "Which version of kubernetes to use. Defaults to skuba latest")
	configCmd.AddCommand(configGenCmd)
}

var configGenCmd = &cobra.Command{
	Use:   "gen [cluster name] [terraform workspace path]",
	Short: "Generate and print the config to the console",
	Long:  ``,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]
		terraformWorkspacePath := args[1]

		if write {
			divingbell.UpdateClusterConfigFile(clusterName, kubernetesVersion, terraformWorkspacePath)
		} else {
			fmt.Println(divingbell.ClusterConfigYamlString(clusterName, kubernetesVersion, terraformWorkspacePath))
		}
	},
}
