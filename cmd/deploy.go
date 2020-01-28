package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"

	divingbell "github.com/tdaines42/diving-bell/pkg/diving-bell"
)

var redeploy bool

func init() {
	deployCmd.Flags().BoolVar(&redeploy, "redeploy", false, "Destroy a cluster and deploy it again")
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

		if redeploy {
			divingbell.DeProvisionCluster(args[1])
		}

		divingbell.ProvisionCluster(args[1])
		divingbell.UpdateClusterConfig(args[0], args[1])

		err := viper.Unmarshal(&config)
		if err != nil {
			klog.Fatalf("unable to decode into struct, %v", err)
		}

		divingbell.BootstrapCluster(config, redeploy)
	},
}
