package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"

	divingbell "github.com/tdaines42/diving-bell/pkg/diving-bell"
)

func init() {
	rootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy [cluster name] [terraform workspace path]",
	Short: "Deploy a cluster with one command",
	Long: `Deploys a cluster with one command using kubectl, terraform, and skuba
	Writes the config to disk`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var config divingbell.ClusterConfig

		divingbell.UpdateClusterConfig(args[0], args[1]) // Write the initial config
		divingbell.ProvisionCluster(viper.GetString("terraformWorkspacePath"))
		divingbell.UpdateClusterConfig(args[0], args[1]) // Update the config with node info

		err := viper.Unmarshal(&config)
		if err != nil {
			klog.Fatalf("unable to decode into struct, %v", err)
		}

		divingbell.BootstrapCluster(config)
	},
}
