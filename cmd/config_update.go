package cmd

import (
	"path"

	"github.com/spf13/cobra"

	divingbell "github.com/tdaines42/diving-bell/pkg/diving-bell"
)

func init() {
	configCmd.AddCommand(configUpdateCmd)
}

var configUpdateCmd = &cobra.Command{
	Use:   "update [cluster name]",
	Short: "Updates the config in the cluster",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig := path.Join(currentWorkingDir, args[0], "admin.conf")
		divingbell.StoreClusterConfig(kubeconfig)
	},
}
