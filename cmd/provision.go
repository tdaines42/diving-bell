package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	divingbell "github.com/tdaines42/diving-bell/pkg/diving-bell"
)

func init() {
	rootCmd.AddCommand(provisionCmd)
}

var provisionCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision the cluster using terraform",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		divingbell.ProvisionCluster(viper.GetString("terraformWorkspacePath"))
	},
}
