package cmd

import (
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
		divingbell.CheckClusterReady(viper.GetString("clusterName"))
	},
}
