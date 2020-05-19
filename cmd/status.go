package cmd

import (
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	divingbell "github.com/tdaines42/diving-bell/pkg/diving-bell"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the current status of the cluster",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig := path.Join(currentWorkingDir, viper.GetString("clusterName"), "admin.conf")
		divingbell.CheckClusterReady(kubeconfig)
	},
}
